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
	if err := c.BindJSON(&body); err != nil {
		log.Err(err).Msg("error parse body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}

	//convert amout from string to float64
	amount, err := strconv.ParseFloat(body.GrossAmount, 64)
	if err != nil {
		log.Err(err).Msg("error convert amount")
		return
	}

	//convert time from string to time.Time
	expiryTime, err := time.Parse("2006-01-02 15:04:05", body.ExpiryTime)
	if err != nil {
		log.Err(err).Msg("error parse expiry time")
	}
	expiryTimeUTC := expiryTime.UTC()

	var settlementTimeUTC *time.Time
	if body.SettlementTime != nil {
		settlementTimeParse, err := time.Parse("2006-01-02 15:04:05", *body.SettlementTime)
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.DBErrorResponse(createPayment.Error)})
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": apiResponse.SuccessSingularResponse(paymentData)})
}
