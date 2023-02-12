package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/trinitt/config"
	"github.com/trinitt/controllers"
	"github.com/trinitt/models"
	"github.com/trinitt/routes"
	sockets "github.com/trinitt/sockets"
	"github.com/trinitt/utils"
)

func main() {

	config.InitConfig()

	config.ConnectDB()
	models.MigrateDB()
	config.CreateProducer()
	config.CreateConsumer()
	server := echo.New()
	hub := sockets.NewHub()
	go hub.Run()
	server.GET("/ws", func(c echo.Context) error {
		sockets.ServeWs(hub, c.Response(), c.Request())
		return nil
	})
	utils.InitLogger(server)
	server.Use(middleware.CORS())
	routes.Init(server)

	controllers.Consume()

	server.Logger.Fatal(server.Start(":" + config.ServerPort))
}
