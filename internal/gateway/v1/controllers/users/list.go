package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/guilherme-de-marchi/coin-commerce/api/users/v1"
	"github.com/guilherme-de-marchi/coin-commerce/internal/gateway/v1/repository/users"
	"github.com/guilherme-de-marchi/coin-commerce/pkg"
)

func (c Controllers) List() {
	c.Group.GET("/list", list)
}

func list(c *gin.Context) {
	var req pb.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, pkg.Error(err, "invalid query params"))
		return
	}

	users.List(c, &req)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
