package services

import (
	"POS-BE/libraries/helpers/api/apiResponse"
	"POS-BE/libraries/helpers/services/midtransService"
	"POS-BE/libraries/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

func CreatePayment(c *gin.Context) {
	//parse body string to object
	body := midtransService.MidtransNotification{}
	log.Info().Any("body", c.Request.Body).Msg("log event")
	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
	}

	//convert amout from string to float64
	amount, err := strconv.ParseFloat(body.GrossAmount, 64)
	if err != nil {
		log.Err(err).Msg("error convert amount")
		return
	}

	// Load zona waktu Jakarta
	loc, _ := time.LoadLocation("Asia/Jakarta")

	//convert time from string to time.Time
	expiryTime, err := time.ParseInLocation("2006-01-02 15:04:05", body.ExpiryTime, loc)
	if err != nil {
		log.Err(err).Msg("error parse expiry time")
	}
	expiryTimeUTC := expiryTime.UTC()

	var settlementTimeUTC *time.Time
	if body.SettlementTime != nil {
		settlementTimeParse, err := time.ParseInLocation("2006-01-02 15:04:05", *body.SettlementTime, loc)
		if err != nil {
			log.Err(err).Msg("error parse settlement time")
		}
		settlementTimeParse = settlementTimeParse.UTC()
		settlementTimeUTC = &settlementTimeParse
	}

	//get acquire, issuer, and reference(va number)
	paymentSource, err := midtransService.DeterminePaymentSource(body)
	if err != nil {
		log.Err(err).Msg("error get payment source")
	}

	//create payment
	paymentData := models.Payment{
		Payment_id:             string(uuid.NewString()),
		Transaction_id:         body.OrderID,
		Payment_method:         body.PaymentType,
		Acquire:                paymentSource.Acquire,
		Issuer:                 paymentSource.Issuer,
		Payment_references:     paymentSource.Payment_references,
		Amount:                 amount,
		Currency:               body.Currency,
		Status:                 body.TransactionStatus,
		Settlement_time:        settlementTimeUTC,
		Expiration_time:        &expiryTimeUTC,
		Gateway_transaction_id: body.TransactionID,
	}

	createPayment := db.Model(&paymentData).Clauses(clause.Returning{}).Create(&paymentData)
	if createPayment.Error != nil {
		log.Err(createPayment.Error).Msg("error create payment")
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(paymentData)})
}
