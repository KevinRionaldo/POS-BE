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

// Struct representing a transaction with additional payment details
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

// Start a new transaction
func StartTransaction(c *gin.Context) {
	// Struct to parse request body
	type bodyType struct {
		Amount      *int64 `json:"amount"`      // Transaction amount
		Description string `json:"description"` // Transaction description
	}

	// Generate a unique transaction ID
	transactionID := string(uuid.NewString())

	// Bind the received JSON to bodyType
	body := bodyType{}
	if err := c.BindJSON(&body); err != nil || body.Amount == nil {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	// Define transaction data
	transactionData := models.Transaction{
		Transaction_id: transactionID,
		Description:    body.Description,
	}

	// Insert transaction record into the database
	createTransaction := db.Model(&transactionData).Clauses(clause.Returning{}).Create(&transactionData)
	if createTransaction.Error != nil {
		log.Err(createTransaction.Error).Msg("error creating transaction")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": apiResponse.GeneralErrorResponse(createTransaction.Error)})
		return
	}

	// Call Midtrans API to initiate payment processing
	midtransResp, midtransErr := midtransService.CreateTransaction(transactionID, *body.Amount)
	if midtransErr != nil {
		log.Err(midtransErr).Msg("midtrans error")
		// If Midtrans fails, delete the transaction from the database
		db.Model(&models.Transaction{}).Where("transaction_id = ?", transactionData.Transaction_id).Delete(&models.Transaction{})
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": apiResponse.GeneralErrorResponse(midtransErr)})
		return
	}

	// Define response data structure
	type responseDataType struct {
		models.Transaction
		GatewayTransactionToken string `json:"gateway_transaction_token"`
		Redirect_URL            string `json:"redirect_URL"`
	}

	// Return successful transaction response
	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessSingularResponse(responseDataType{
		Transaction:             transactionData,
		GatewayTransactionToken: midtransResp.Token,
		Redirect_URL:            midtransResp.RedirectURL,
	})})
}

// Retrieve a paginated list of transactions
func GetTransactions(c *gin.Context) {
	// Set pagination parameters
	page, limit := paging.SetPageLimit(c.Query("page"), c.Query("limit"))

	listTransactions := []transactionType{}

	// SQL query to fetch transactions with payment details
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

	// Create subquery to get transactions with payment details
	getTransactionSubQuery := db.Select(selectColumns).
		Table(config.GetTableNameOnCurrentSchema(`transaction as t`)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON p.transaction_id = t.transaction_id", config.GetTableNameOnCurrentSchema("payment p"))).
		Order("t.transaction_id, t.created_at desc, p.created_at DESC")

	// Final query to retrieve transactions with pagination
	getTransaction := db.Select("*").
		Table("(?) as subquery", getTransactionSubQuery).
		Order("created_at DESC, transaction_id asc").
		Limit(limit).Offset((page - 1) * limit).
		Find(&listTransactions)

	if getTransaction.Error != nil {
		log.Err(getTransaction.Error).Msg("error retrieving transactions from PostgreSQL")
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": apiResponse.DBErrorResponse(getTransaction.Error)})
		return
	}

	// Count total transactions
	var totalTransaction int64
	countTransaction := db.Select(`COUNT(DISTINCT t.transaction_id)`).
		Table(config.GetTableNameOnCurrentSchema(`transaction as t`)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON p.transaction_id = t.transaction_id", config.GetTableNameOnCurrentSchema("payment p"))).
		Find(&totalTransaction)

	if countTransaction.Error != nil {
		log.Err(getTransaction.Error).Msg("error counting transactions in PostgreSQL")
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": apiResponse.DBErrorResponse(countTransaction.Error)})
		return
	}

	// Return paginated transaction list response
	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessPluralResponse(listTransactions, totalTransaction, limit, page)})
}

// Retrieve transaction details by transaction ID
func GetTransactionsByID(c *gin.Context) {
	// Define query filter to get only settled payments
	qryFilter := `t.transaction_id = $1 AND p.status = 'settlement'`
	qryValue := []any{c.Param("id")}

	transactionData := transactionType{}

	// SQL query to fetch transaction and payment details
	selectColumns := ` t.*,
    p.payment_method,
    p.acquire,
    p.payment_references,
    p.amount,
    p.currency,
    p."status",
    p.settlement_time`

	// Execute query
	getTransaction := db.Select(selectColumns).
		Table(config.GetTableNameOnCurrentSchema(`transaction as t`)).
		Joins(fmt.Sprintf("LEFT JOIN %s ON p.transaction_id = t.transaction_id", config.GetTableNameOnCurrentSchema("payment p"))).
		Where(qryFilter, qryValue...).
		Order("p.created_at DESC").
		Find(&transactionData)

	if getTransaction.Error != nil {
		log.Err(getTransaction.Error).Msg("error retrieving transaction details from PostgreSQL")
		c.IndentedJSON(http.StatusBadGateway, gin.H{"message": apiResponse.DBErrorResponse(getTransaction.Error)})
		return
	}

	// Return transaction details
	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessSingularResponse(transactionData)})
}
