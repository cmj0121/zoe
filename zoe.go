package zoe

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// The project name of the Zoe.
	PROJ_NAME = "zoe"

	// The version of the Zoe.
	MAJOR = 0
	MINOR = 2
	MICRO = 0
)

// The Zoe instance that holds the CLI and the logger.
type Zoe struct {
	Version kong.VersionFlag `short:"V" help:"Show version information and exit"`

	// The logger options.
	Verbose int  `short:"v" xor:"quite,verbose" type:"counter" help:"Show the verbose output" default:"0"`
	Quiet   bool `short:"q" xor:"quite,verbose" help:"Show no output"`
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

	return z.run()
}

func (z *Zoe) run() error {
	return nil
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
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}

func (z *Zoe) epilogue() {
}
