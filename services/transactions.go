package services

import (
	"POS-BE/libraries/config"
	"POS-BE/libraries/helpers/api/apiResponse"
	"POS-BE/libraries/helpers/services/midtransService"
	"POS-BE/libraries/helpers/utils/paging"
	"POS-BE/libraries/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

type transactionType struct {
	models.Transaction
	Payment_method     string     `gorm:"not null" json:"payment_method"`
	Acquire            string     `json:"acquire"`
	Payment_references *string    `json:"payment_references"`
	Amount             float64    `gorm:"type:numeric(15,2);not null" json:"amount"`
	Currency           string     `gorm:"not null;default:'IDR'" json:"currency"`
	Status             string     `json:"status"`
	Settlement_time    *time.Time `json:"settlement_time"`
}

func StartTransaction(c *gin.Context) {
	type bodyType struct {
		Amount      *int64 `json:"amount"`
		Description string `json:"description"`
	}

	transactionID := string(uuid.NewString())
	log.Info().Any("transactionID", transactionID).Msg("log transaction ID")
	// Call BindJSON to bind the received JSON to
	// body data.
	body := bodyType{}
	if err := c.BindJSON(&body); err != nil || body.Amount == nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	//declare transaction data
	transactionData := models.Transaction{
		Transaction_id: transactionID,
		Description:    body.Description,
	}

	createTransaction := db.Model(&transactionData).Clauses(clause.Returning{}).Create(&transactionData)
	if createTransaction.Error != nil {
		log.Err(createTransaction.Error).Msg("error create transaction")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": apiResponse.GeneralErrorResponse(createTransaction.Error)})
		return
	}

	midtransResp, midtransErr := midtransService.CreateTransaction(transactionID, *body.Amount)
	if midtransErr != nil {
		log.Err(midtransErr).Msg("midtrans error")
		db.Model(&models.Transaction{}).Where("transaction_id = ?", transactionData.Transaction_id).Delete(&models.Transaction{})
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": apiResponse.GeneralErrorResponse(midtransErr)})
		return
	}

	type responseDataType struct {
		models.Transaction
		GatewayTransactionToken string `json:"gateway_transaction_token"`
		Redirect_URL            string `json:"redirect_URL"`
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessSingularResponse(responseDataType{
		Transaction:             transactionData,
		GatewayTransactionToken: midtransResp.Token,
		Redirect_URL:            midtransResp.RedirectURL,
	})})
}

func GetTransactions(c *gin.Context) {
	//pagination
	page, limit := paging.SetPageLimit(c.Query("page"), c.Query("limit"))

	listTransactions := []transactionType{}

	//query builder for get transaction, and payment
	selectColumns := `DISTINCT ON(t.transaction_id) 
    t.*,
    p.payment_method,
    p.acquire,
    p.payment_references,
    p.amount,
    p.currency,
	CASE 
        WHEN p.status IS NULL THEN 'payment_selection'
        ELSE p."status"  
    END AS "status", 
    p.settlement_time`
	getTransactionSubQuery := db.Select(selectColumns).
		Table(config.GetTableNameOnCurrentSchema(`transaction as t`)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON p.transaction_id = t.transaction_id", config.GetTableNameOnCurrentSchema("payment p"))).
		Order("t.transaction_id, t.created_at desc, p.created_at DESC")

	getTransaction := db.Select("*").
		Table("(?) as subquery", getTransactionSubQuery).
		Order("created_at DESC, transaction_id asc").
		Limit(limit).Offset((page - 1) * limit).
		Find(&listTransactions)

	if getTransaction.Error != nil {
		log.Err(getTransaction.Error).Msg("error get transaction postgre")
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": apiResponse.DBErrorResponse(getTransaction.Error)})
		return
	}

	var totalTransaction int64
	countTransaction := db.Select(`COUNT(DISTINCT t.transaction_id)`).
		Table(config.GetTableNameOnCurrentSchema(`transaction as t`)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON p.transaction_id = t.transaction_id", config.GetTableNameOnCurrentSchema("payment p"))).
		Find(&totalTransaction)

	if countTransaction.Error != nil {
		log.Err(getTransaction.Error).Msg("error count transaction postgre")
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": apiResponse.DBErrorResponse(countTransaction.Error)})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessPluralResponse(listTransactions, totalTransaction, limit, page)})
}

func GetTransactionsByID(c *gin.Context) {
	qryFilter := `t.transaction_id = $1 AND p.status = 'settlement'`
	qryValue := []any{c.Param("id")}

	transactionData := transactionType{}

	//query builder for get transaction, payment, and generated token
	selectColumns := ` t.*,
    p.payment_method,
    p.acquire,
    p.payment_references,
    p.amount,
    p.currency,
    p."status",
    p.settlement_time`
	getTransaction := db.Select(selectColumns).
		Table(config.GetTableNameOnCurrentSchema(`transaction as t`)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON p.transaction_id = t.transaction_id", config.GetTableNameOnCurrentSchema("payment p"))).
		Where(qryFilter, qryValue...).
		Order("p.created_at DESC").
		Find(&transactionData)

	if getTransaction.Error != nil {
		log.Err(getTransaction.Error).Msg("error get transaction postgre")
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": apiResponse.DBErrorResponse(getTransaction.Error)})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(transactionData)})
}
