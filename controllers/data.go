package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/trinitt/config"
	"github.com/trinitt/models"
	"github.com/trinitt/utils"
)

type AddDatasetRequest struct {
	EntityID uint                   `json:"entityId"`
	Dataset  []map[uint]interface{} `json:"dataset"`
}

func AddDataSet(c echo.Context) error {
	var req AddDatasetRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Error parsing request body")
	}

	db := config.GetDB()

	userId := utils.GetUserID(c)

	var parameters []models.Parameter

	var maxEntriesFloat float64

	db.Model(&models.Data{}).Where("entity_id = ?", req.EntityID).Select("MAX(row)").Row().Scan(&maxEntriesFloat)

	maxEntries := uint(maxEntriesFloat) + 1

	if err := db.Preload("Entity").Where("entity_id = ?", req.EntityID).Find(&parameters).Error; err != nil {
		return c.JSON(http.StatusBadRequest, "Error finding parameters for this entity")
	}

	if len(parameters) == 0 {
		return c.JSON(http.StatusBadRequest, "No parameters found for this entity")
	}

	if parameters[0].Entity.CreatedByID != userId {
		return c.JSON(http.StatusUnauthorized, "You are not authorized to add data to this entity")
	}

	var dataset []models.Data
	var records []Record

	for i, data := range req.Dataset {
		var record Record
		record.User_id = int(userId)
		record.Entity_id = int(req.EntityID)

		for _, parameter := range parameters {
			paramType := parameter.Type
			if paramType == models.ParameterTypeString {
				if val, ok := data[parameter.ID].(string); ok {
					dataset = append(dataset, models.Data{
						Value:       val,
						ParameterID: parameter.ID,
						EntityID:    parameter.EntityID,
						Row:         maxEntries + uint(i),
					})

					newParam := ParamType{
						Data_type: "string",
						Value:     val,
					}

					record.Param = append(record.Param, newParam)

				}

			} else if paramType == models.ParameterTypeInt {
				val := strconv.Itoa(int(data[parameter.ID].(float64)))
				dataset = append(dataset, models.Data{
					Value:       val,
					ParameterID: parameter.ID,
					EntityID:    parameter.EntityID,
					Row:         maxEntries + uint(i),
				})
				newParam := ParamType{
					Data_type: "int",
					Value:     val,
				}
				record.Param = append(record.Param, newParam)
			}
		}

		records = append(records, record)
	}

	go Produce(records)

	if err := db.Create(&dataset).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Error adding data")
	}

	return c.JSON(http.StatusOK, "Data added successfully")
}
