package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinitt/config"
	"github.com/trinitt/models"
	"github.com/trinitt/utils"
)

type CreateEntityRequest struct {
	Name string `json:"name"`
}

type EntityResponse struct {
	Name string `json:"name" gorm:"column:name"`
	ID   uint   `json:"id" gorm:"column:id"`
}

func CreateEntity(c echo.Context) error {
	var req CreateEntityRequest

	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, "Error-1")
	}

	db := config.GetDB()

	userId := utils.GetUserID(c)

	entity := models.Entity{
		Name:        req.Name,
		CreatedByID: userId,
	}

	if err := db.Create(&entity).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Error-2")

	}

	var response EntityResponse

	response.Name = entity.Name
	response.ID = entity.ID

	return c.JSON(http.StatusOK, response)
}

func GetEntities(c echo.Context) error {
	userId := utils.GetUserID(c)

	db := config.GetDB()

	var entities []EntityResponse

	if err := db.Model(models.Entity{}).Where("created_by_id = ?", userId).Find(&entities).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Error-1")
	}

	return c.JSON(http.StatusOK, entities)
}
