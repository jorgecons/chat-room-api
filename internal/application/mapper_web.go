package application

import (
	"github.com/gin-gonic/gin"
)

func (a *App) MapWebRoutes() *App {
	a.router.StaticFile("/static", "./static")
	a.router.POST("/users", a.handlers.CreateAccountHandler)
	a.router.POST("/login", a.handlers.LoginHandler)
	a.router.GET("", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	return a
}
