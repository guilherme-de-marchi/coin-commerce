package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
	"github.com/rs/zerolog/log"
)

func HandleErrors(c *gin.Context) {
	c.Next()

	var errors []pkg.Err
	var out []gin.H
	for _, err := range c.Errors {
		var msg string
		var errx pkg.Err
		switch e := err.Err.(type) {
		case pkg.Err:
			msg = e.PublicMsg
			errx = e
		default:
			msg = "unknown error, contact support"
			errx = pkg.Error(err.Err).(pkg.Err)
		}
		errors = append(errors, errx)
		out = append(out, gin.H{"error": msg})
	}

	if len(errors) != 0 {
		log.Error().Any("errors", errors).Msg("")
		c.JSON(-1, out)
	}
}
