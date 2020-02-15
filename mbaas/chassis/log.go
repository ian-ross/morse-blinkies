package chassis

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogSetup performs logging setup common to all services.
func LogSetup(dev bool) {
	baselog := zerolog.New(os.Stdout)
	if dev {
		baselog = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	applog := baselog.With().Timestamp().Logger()
	log.Logger = applog
}
