package zoe

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/cmj0121/zoe/pkg/honeypot"
	"github.com/cmj0121/zoe/pkg/monitor"
)

const (
	// The project name of the Zoe.
	PROJ_NAME = "zoe"

	// The version of the Zoe.
	MAJOR = 0
	MINOR = 2
	MICRO = 9
)

// The Zoe instance that holds the CLI and the logger.
type Zoe struct {
	Version kong.VersionFlag `short:"V" help:"Show version information and exit"`

	// The logger options.
	Verbose int  `short:"v" xor:"quite,verbose" type:"counter" help:"Show the verbose output" default:"0"`
	Quiet   bool `short:"q" xor:"quite,verbose" help:"Show no output"`

	// The external configuration
	Config   *string         `short:"c" help:"The external configuration file"`
	Database *Database       `embed:"" help:"The database service"`
	Server   *monitor.Server `embed:"" help:"The MongoDB service"`

	Service honeypot.Service `arg:"" help:"The honeypot service" default:"ssh"`
}

func init() {
	// setup the default logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	writer := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = zerolog.New(writer).With().Timestamp().Logger()
}

// New creates a new Zoe instance with the default logger.
func New() *Zoe {
	return &Zoe{}
}

// Parse the command line arguments and run the command.
func (z *Zoe) ParseAndRun() error {
	opts := []kong.Option{
		kong.Name("zeo"),
		kong.Description("The simple but all-in-one honeypot service."),
		kong.Vars{"version": fmt.Sprintf("%s/%d.%d.%d", PROJ_NAME, MAJOR, MINOR, MICRO)},
	}

	kong.Parse(z, opts...)
	return z.Run()
}

// Run the Zoe instance with the known arguments.
func (z *Zoe) Run() error {
	z.prologue()
	defer z.epilogue()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// starting the graceful shutdown, catch the signal SIGINT and SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		defer signal.Stop(sig)
		defer close(sig)

		select {
		case <-sig:
			log.Info().Msg("received the signal, starting the graceful shutdown ...")
			cancel()
		case <-ctx.Done():
		}
	}()

	go z.Server.Run(ctx)
	return z.Service.Run(ctx)
}

func (z *Zoe) prologue() {
	if z.Quiet {
		zerolog.SetGlobalLevel(zerolog.Disabled)
		return
	}

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

	log.Info().Msg("finished the prologue ...")
	z.loadConfig()
	z.Database.Init()
	z.Database.Migrate()
}

func (z *Zoe) epilogue() {
	log.Info().Msg("starting the epilogue ...")
	log.Info().Msg("finished the epilogue ...")
}

func (z *Zoe) loadConfig() {
	if z.Config == nil {
		log.Debug().Msg("skip the external configuration")
		return
	}

	v := viper.New()
	v.SetConfigFile(*z.Config)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Warn().Err(err).Msg("failed to read the external configuration")
		return
	}

	// override the configuration by the external configuration
	if err := v.Unmarshal(z); err != nil {
		log.Warn().Err(err).Msg("failed to unmarshal the external configuration")
		return
	}

	// override the configuration by the external configuration
	z.Service.Viper = v.Sub("service")
}
