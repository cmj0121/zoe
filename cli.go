package zoe

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	kong.Parse(z, options...)
	return z.Run()
}

// Run the Zoe with the given configuration
func (z *Zoe) Run() int {
	z.prologue()
	defer z.epilogue()

	return z.run()
}

// run the zoe service which already setup everything well.
func (z *Zoe) run() int {
	log.Info().Msg("starting zoe ...")
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
