package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/guilherme-de-marchi/coin-commerce/internal/gateway/v1/controllers/orders"
	"github.com/guilherme-de-marchi/coin-commerce/internal/gateway/v1/controllers/users"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
)

func Setup(g *gin.RouterGroup) {
	pkg.SetupControllers(users.Controllers{Group: g.Group("/users")})
	pkg.SetupControllers(orders.Controllers{Group: g.Group("/orders")})
}
