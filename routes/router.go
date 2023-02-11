package routes

import (
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinitt/controllers"
	"github.com/trinitt/sockets"
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

func generateRandomColor() string {
	colors := []string{"red", "green", "blue"}
	return colors[rand.Intn(len(colors))]
}

func generateRandomNodes() []Node {
	nodes := []Node{}
	for i := 0; i < 100; i++ {
		nodes = append(nodes, Node{
			ID:    i,
			X:     rand.Intn(1000),
			Y:     rand.Intn(1000),
			Color: generateRandomColor(),
		})
	}
	return nodes
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
	api.GET("/kafka", controllers.Produce)
	api.GET("/nodes", func(c echo.Context) error {
		nodes := []Node{
			{ID: 1, X: 1, Y: 1, Color: "red"},
			{ID: 2, X: 2, Y: 2, Color: "green"},
			{ID: 3, X: 3, Y: 3, Color: "blue"},
		}
		return c.JSON(200, nodes)
	})
	api.POST("/nodes", func(c echo.Context) error {
		var req SendNodesRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		nodes := generateRandomNodes()

		go SendStreamOfNodes(req.UserID, nodes)

		return c.JSON(200, "OK")
	})

	api.POST("/entity/create", controllers.CreateEntity)
	api.POST("/parameter/create", controllers.AddParamToEntity)
	api.POST("/data/add", controllers.AddData)

}
