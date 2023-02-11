package controllers

import "github.com/labstack/echo/v4"

type Node struct {
	ID    int    `json:"id"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
}

func GetNodesByEntityID(c echo.Context) error {
	// entityId := c.Param("entityId")

	// return c.JSON(200, nodes)
	return nil
}
