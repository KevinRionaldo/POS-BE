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

func CreateTransactionProduct(c *gin.Context) {
	//parse transactionProductData string to object
	transactionProductData := models.Transaction_product{}
	if err := c.BindJSON(&transactionProductData); err != nil {
		log.Err(err).Msg("error parse body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}
	transactionProductData.Transaction_product_id = string(uuid.NewString())

	createTransactionProduct := db.Model(&transactionProductData).Clauses(clause.Returning{}).Create(&transactionProductData)
	if createTransactionProduct.Error != nil {
		log.Err(createTransactionProduct.Error).Msg("error create transactionProduct")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(createTransactionProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(transactionProductData)})
}

func GetTransactionProduct(c *gin.Context) {
	//get product by categories id filter
	qryFilter := ""
	qryValue := []any{}

	//params filter for transaction id and product id
	if c.Query("transaction_id") != "" {
		qryValue = append(qryValue, c.Query("transaction_id"))
		qryFilter += fmt.Sprintf("transaction_id = $%d", len(qryValue))
	}
	if c.Query("product_id") != "" {
		qryValue = append(qryValue, c.Query("product_id"))
		qryFilter += fmt.Sprintf("product_id = $%d", len(qryValue))
	}

	//pagination
	page, limit := paging.SetPageLimit(c.Query("page"), c.Query("limit"))

	//parse transactionProductData string to object
	transactionProductData := []models.Transaction_product{}

	//build query for getTransactionProduct
	getTransactionProduct := db.Model(&transactionProductData).
		Where(qryFilter, qryValue...).
		Order("created_at ASC").
		Limit(limit).Offset((page - 1) * limit).
		Find(&transactionProductData)

	if getTransactionProduct.Error != nil {
		log.Err(getTransactionProduct.Error).Msg("error get transactionProduct")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(getTransactionProduct.Error)})
		return
	}

	//build query for totalTransactionProduct
	var totalTransactionProduct int64
	countTransactionProduct := db.Model(&transactionProductData).Where(qryFilter, qryValue...).Count(&totalTransactionProduct)
	if getTransactionProduct.Error != nil {
		log.Err(getTransactionProduct.Error).Msg("error get count transactionProduct")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(countTransactionProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessPluralResponse(transactionProductData, totalTransactionProduct, limit, page)})
}

func UpdateTransactionProduct(c *gin.Context) {
	//parse transactionProductData string to object
	transactionProductData := models.Transaction_product{}

	if err := c.BindJSON(&transactionProductData); err != nil {
		log.Err(err).Msg("error parse body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}
	transactionProductData.Transaction_product_id = c.Param("id")

	//build query for updateTransactionProduct
	updateTransactionProduct := db.Model(&transactionProductData).Clauses(clause.Returning{}).Updates(&transactionProductData)
	if updateTransactionProduct.Error != nil {
		log.Err(updateTransactionProduct.Error).Msg("error update transactionProduct")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(updateTransactionProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessSingularResponse(transactionProductData)})
}

func DeleteTransactionProduct(c *gin.Context) {
	//parse transactionProductData string to object
	transactionProductData := models.Transaction_product{Transaction_product_id: c.Param("id")}

	//build query for deleteTransactionProduct
	deleteTransactionProduct := db.Model(&transactionProductData).Clauses(clause.Returning{}).Delete(&transactionProductData)
	if deleteTransactionProduct.Error != nil {
		log.Err(deleteTransactionProduct.Error).Msg("error delete transactionProduct")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(deleteTransactionProduct.Error)})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessSingularResponse(transactionProductData)})
}
