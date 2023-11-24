package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/guilherme-de-marchi/coin-commerce/internal/gateway/middlewares"
	"github.com/guilherme-de-marchi/coin-commerce/internal/gateway/v1/controllers"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func Start() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	pkg.Globals.RabbitMQ = conn

	if err = pkg.Globals.Validate(); err != nil {
		panic(err)
	}

	e := gin.New()
	e.Use(
		gin.Recovery(),
		middlewares.HandleErrors,
	)
	controllers.Setup(e.Group("/v1"))
	e.Run(":8080")
}
