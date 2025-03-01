package application

func (a *App) MapWebRoutes() *App {
	a.router.POST("/users", a.handlers.CreateAccountHandler)
	a.router.POST("/login", a.handlers.LoginHandler)
	return a
}
