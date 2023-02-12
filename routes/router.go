package routes

import (
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/trinitt/controllers"
	"github.com/trinitt/sockets"
	"github.com/trinitt/utils"
)

type Node struct {
	ID    int    `json:"id"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
}

type SendNodesRequest struct {
	UserID uint `json:"userId"`
}

func SendStreamOfNodes(userID uint, nodes []Node) {
	for _, v := range nodes {
		message := sockets.Message{
			Type: "NODE",
			Data: map[string]interface{}{
				"node": v,
			},
		}

		sockets.SendMessageToClient(userID, message.Type, message.Data)

		time.Sleep(1 * time.Second)
	}
}

func Init(e *echo.Echo) {
	api := e.Group("/api")

	UserRoutes(api)
	api.Use(echojwt.WithConfig(utils.JWTConfig))

	api.POST("/entity/create", controllers.CreateEntity)
	api.POST("/parameter/create", controllers.AddParamToEntity)

	api.GET("/entities", controllers.GetEntities)
	api.POST("/dataset/add", controllers.AddDataSet)
	api.GET("/params/:entityId", controllers.GetParametersByEntityID)
	api.GET("/nodes/:entityId", controllers.GetNodesByEntityID)

}
