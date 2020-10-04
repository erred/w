package usvc

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/rs/zerolog"
)

type LoggerOpts struct {
	LogLevel  string
	LogFormat string
}

func (o *LoggerOpts) Flags(fs *flag.FlagSet) {
	fs.StringVar(&o.LogLevel, "log.lvl", "trace", "log level: trace, debug, info, error")
	fs.StringVar(&o.LogFormat, "log.fmt", "json", "log format: logfmt, json")
}

// Logger creates a logger and optionally sets the global logger
func (o LoggerOpts) Logger(global bool) zerolog.Logger {
	var logout io.Writer = os.Stderr
	lvl, _ := zerolog.ParseLevel(o.LogLevel)
	if o.LogFormat == "logfmt" {
		logout = zerolog.ConsoleWriter{Out: logout}
	}

	lg := zerolog.New(logout).Level(lvl).With().Timestamp().Logger()
	if global {
		log.SetOutput(lg)
	}
	return lg
}
