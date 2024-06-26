package zoe

import (
	"fmt"
	"os"
	"sync"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zoe/pkg/config"
	"github.com/cmj0121/zoe/pkg/service/types"
)

const (
	PROJ_NAME = "zoe"
)

const (
	MAJOR = 0
	MINOR = 1
	MAINT = 0
)

// The zoe instance is the main entry point for the Zoe CLI
type Zoe struct {
	// general options
	Version kong.VersionFlag `short:"V" name:"version" help:"Show version information and exit."`

	// logging options
	Verbose int `short:"v" name:"verbose" type:"counter" help:"Show the verbose output."`

	Config config.Config `short:"c" name:"config" help:"The configuration file of the zoe service."`

	// execute migration
	Migrate struct {
		// the target database to execute
		Database string `arg:"" help:"The target database to execute."`
	} `cmd:"" help:"Execute the migration to the latest version."`

	// launch the zoe service
	Launch struct{} `cmd:"" help:"Launch the zoe service."`

	// Show the current configuration
	Show struct{} `cmd:"" help:"Show the current configuration."`
}

func init() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
}

// Create a new Zoe instance with default configuration
func New() *Zoe {
	return &Zoe{}
}

// Parse the command line arguments and run the Zoe CLI
func (z *Zoe) ParseAndRun() int {
	options := []kong.Option{
		kong.Name("zoe"),
		kong.Description("The simple but all-in-one honeypot service."),
		kong.Vars{
			"version": fmt.Sprintf("%s (v%d.%d.%d)", PROJ_NAME, MAJOR, MINOR, MAINT),
		},
	}

	ctx := kong.Parse(z, options...)
	return z.Run(ctx.Command())
}

// Run the Zoe with the given configuration
func (z *Zoe) Run(command string) int {
	z.prologue()
	defer z.epilogue()

	switch command {
	case "migrate <database>":
		return MigrateUp(z.Migrate.Database)
	case "show":
		fmt.Println(z.Config.ToYAML())
		return 0
	default:
		return z.run()
	}
}

// run the zoe service which already setup everything well.
func (z *Zoe) run() int {
	log.Info().Msg("starting zoe ...")

	community := make(chan *types.Message)
	defer close(community)

	// start receiving the message from the community channel
	go z.Config.Loggers.Run(community)

	// start the monitor service
	if z.Config.Monitor != nil {
		log.Info().Str("bind", z.Config.Monitor.Bind).Msg("starting the monitor ...")
		go z.Config.Monitor.Serve()
	}

	wg := sync.WaitGroup{}
	for svc_name, svc := range z.Config.Services {
		log.Info().Str("service", svc_name).Msg("running the service ...")

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := svc.Run(community); err != nil {
				log.Error().Err(err).Str("service", svc_name).Msg("failed to run the service")
				return
			}
		}()
	}

	wg.Wait()
	return 0
}

// setup everything before running the zoe service
func (z *Zoe) prologue() {
	// setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	writer := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()

	switch z.Verbose {
	case 0:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	log.Debug().Msg("finished set up zoe ...")
}

// clean-up everything after running the zoe service
func (z *Zoe) epilogue() {
	log.Debug().Msg("starting clean up zoe ...")
	log.Debug().Msg("finished clean up zoe ...")
}
