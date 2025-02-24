package services

import (
	"POS-BE/libraries/helpers/api/apiResponse"
	"POS-BE/libraries/helpers/utils/paging"
	"POS-BE/libraries/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

func CreateProducts(c *gin.Context) {
	//parse productData string to object
	productData := models.Product{}
	if err := c.BindJSON(&productData); err != nil || productData.Categories_id == nil {
		log.Err(err).Msg("error parse body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}
	productData.Product_id = string(uuid.NewString())

	createProduct := db.Model(&productData).Clauses(clause.Returning{}).Create(&productData)
	if createProduct.Error != nil {
		log.Err(createProduct.Error).Msg("error create product")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(createProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(productData)})
}

func GetProducts(c *gin.Context) {
	//get product by categories id filter
	qryFilter := ""
	qryValue := []any{}
	if c.Query("categories_id") != "" {
		qryValue = append(qryValue, c.Query("categories_id"))
		qryFilter += fmt.Sprintf("categories_id = $%d", len(qryValue))
	}

	//pagination
	page, limit := paging.SetPageLimit(c.Query("page"), c.Query("limit"))

	//declare productData
	productData := []models.Product{}

	//build query for getProduct
	getProduct := db.Model(&productData).
		Where(qryFilter, qryValue...).
		Order("categories_id ASC, created_at ASC").
		Limit(limit).Offset((page - 1) * limit).
		Find(&productData)

	if getProduct.Error != nil {
		log.Err(getProduct.Error).Msg("error get product")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(getProduct.Error)})
		return
	}

	//build query for totalProduct
	var totalProduct int64
	countProduct := db.Model(&productData).Where(qryFilter, qryValue...).Count(&totalProduct)
	if getProduct.Error != nil {
		log.Err(getProduct.Error).Msg("error get count product")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(countProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessPluralResponse(productData, totalProduct, limit, page)})
}

func UpdateProducts(c *gin.Context) {
	//parse productData string to object
	productData := models.Product{}

	if err := c.BindJSON(&productData); err != nil {
		log.Err(err).Msg("error parse body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}
	productData.Product_id = c.Param("id")

	//build query for updateProduct
	updateProduct := db.Model(&productData).Clauses(clause.Returning{}).Updates(&productData)
	if updateProduct.Error != nil {
		log.Err(updateProduct.Error).Msg("error update product")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(updateProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(productData)})
}

func DeleteProducts(c *gin.Context) {
	//declare productData
	productData := models.Product{Product_id: c.Param("id")}

	//build query for deleteProduct
	deleteProduct := db.Model(&productData).Clauses(clause.Returning{}).Delete(&productData)
	if deleteProduct.Error != nil {
		log.Err(deleteProduct.Error).Msg("error delete product")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(deleteProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(productData)})
}
