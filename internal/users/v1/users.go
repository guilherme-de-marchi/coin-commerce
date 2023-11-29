package v1

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Start() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	s := NewServer()
	mb := NewMessageBroker()

	if err := mb.Run(s); err != nil {
		log.Panic().Err(err).Msg("")
	}
}
