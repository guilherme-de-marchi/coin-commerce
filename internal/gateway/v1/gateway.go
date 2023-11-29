package gateway

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-de-marchi/coin-commerce/internal/gateway/v1/controllers"
	"github.com/guilherme-de-marchi/coin-commerce/internal/gateway/v1/middlewares"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Start() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	mb, err := pkg.NewMessageBroker("gateway")
	if err != nil {
		log.Panic().Err(err).Msg("")
		return
	}
	defer mb.Conn.Close()
	go mb.Start(context.Background(), 5)
	pkg.Globals.MessageBroker = mb

	if err = pkg.Globals.Validate(); err != nil {
		log.Panic().Err(err).Msg("")
		return
	}

	e := gin.New()
	e.Use(
		gin.Recovery(),
		middlewares.HandleErrors,
	)
	controllers.Setup(e.Group("/v1"))

	if err = e.Run(":8080"); err != nil {
		log.Panic().Err(err).Msg("")
		return
	}
}
