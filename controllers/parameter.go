package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinitt/config"
	"github.com/trinitt/models"
	"github.com/trinitt/utils"
)

type AddParamToEntityRequest struct {
	Name     string               `json:"name"`
	Type     models.ParameterType `json:"type"`
	EntityID uint                 `json:"entityId"`
}

type AddParamToEntityResponse struct {
	Name     string               `json:"name"`
	Type     models.ParameterType `json:"type"`
	EntityID uint                 `json:"entityId"`
}

type ParamResponse struct {
	ID      uint                 `json:"id"`
	KeyName string               `json:"keyName"`
	Type    models.ParameterType `json:"type"`
}

func AddParamToEntity(c echo.Context) error {
	var req AddParamToEntityRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db := config.GetDB()

	userId := utils.GetUserID(c)

	var entity models.Entity

	if err := db.Where("id = ?", req.EntityID).First(&entity).Error; err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if entity.CreatedByID != userId {
		return c.JSON(http.StatusUnauthorized, "You are not authorized to add parameters to this entity")
	}

	parameter := models.Parameter{
		KeyName:  req.Name,
		Type:     req.Type,
		EntityID: entity.ID,
	}

	if err := db.Create(&parameter).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	res := AddParamToEntityResponse{
		Name:     parameter.KeyName,
		Type:     parameter.Type,
		EntityID: parameter.EntityID,
	}

	return c.JSON(http.StatusOK, res)
}

func GetParametersByEntityID(c echo.Context) error {
	entityId := c.Param("entityId")

	db := config.GetDB()

	var res []ParamResponse

	if err := db.Model(models.Parameter{}).Where("entity_id = ?", entityId).Find(&res).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}
