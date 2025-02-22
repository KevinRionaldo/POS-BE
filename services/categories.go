package services

import (
	"POS-BE/libraries/helpers/api/apiResponse"
	"POS-BE/libraries/helpers/utils/paging"
	"POS-BE/libraries/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

func CreateCategories(c *gin.Context) {
	//parse categoriesData string to object
	categoriesData := models.Categories{}
	if err := c.BindJSON(&categoriesData); err != nil {
		log.Err(err).Msg("error parse body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}
	categoriesData.Categories_id = string(uuid.NewString())

	createCategories := db.Model(&categoriesData).Clauses(clause.Returning{}).Create(&categoriesData)
	if createCategories.Error != nil {
		log.Err(createCategories.Error).Msg("error create categories")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(createCategories.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(categoriesData)})
}

func GetCategories(c *gin.Context) {
	//pagination
	page, limit := paging.SetPageLimit(c.Query("page"), c.Query("limit"))

	//parse categoriesData string to object
	categoriesData := []models.Categories{}

	//build query for getCategories
	getCategories := db.Model(&categoriesData).
		Order("created_at ASC").
		Limit(limit).Offset((page - 1) * limit).
		Find(&categoriesData)

	if getCategories.Error != nil {
		log.Err(getCategories.Error).Msg("error get categories")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(getCategories.Error)})
		return
	}

	//build query for totalCategories
	var totalCategories int64
	countCategories := db.Model(&categoriesData).Count(&totalCategories)
	if getCategories.Error != nil {
		log.Err(getCategories.Error).Msg("error get count categories")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(countCategories.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessPluralResponse(categoriesData, totalCategories, limit, page)})
}

func UpdateCategories(c *gin.Context) {
	//parse categoriesData string to object
	categoriesData := models.Categories{}

	if err := c.BindJSON(&categoriesData); err != nil {
		log.Err(err).Msg("error parse body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}
	categoriesData.Categories_id = c.Param("id")

	//build query for updateCategories
	updateCategories := db.Model(&categoriesData).Clauses(clause.Returning{}).Updates(&categoriesData)
	if updateCategories.Error != nil {
		log.Err(updateCategories.Error).Msg("error update categories")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(updateCategories.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(categoriesData)})
}

func DeleteCategories(c *gin.Context) {
	//parse categoriesData string to object
	categoriesData := models.Categories{Categories_id: c.Param("id")}

	//build query for deleteCategories
	deleteCategories := db.Model(&categoriesData).Clauses(clause.Returning{}).Delete(&categoriesData)
	if deleteCategories.Error != nil {
		log.Err(deleteCategories.Error).Msg("error delete categories")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(deleteCategories.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(categoriesData)})
}
