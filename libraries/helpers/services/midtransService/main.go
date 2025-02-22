package midtransService

import (
	"fmt"

	"POS-BE/libraries/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransNotification struct {
	BillKey           *string             `json:"bill_key"`
	BillCode          *string             `json:"bill_code"`
	Acquirer          string              `json:"acquirer"`
	Issuer            *string             `json:"issuer"`
	Currency          string              `json:"currency"`
	ExpiryTime        string              `json:"expiry_time"`
	FraudStatus       string              `json:"fraud_status"`
	GrossAmount       string              `json:"gross_amount"`
	MerchantID        string              `json:"merchant_id"`
	ReferenceID       *string             `json:"reference_id"`
	OrderID           string              `json:"order_id"`
	PaymentType       string              `json:"payment_type"`
	SettlementTime    *string             `json:"settlement_time"`
	SignatureKey      string              `json:"signature_key"`
	StatusCode        string              `json:"status_code"`
	StatusMessage     string              `json:"status_message"`
	TransactionID     string              `json:"transaction_id"`
	TransactionStatus string              `json:"transaction_status"`
	TransactionTime   string              `json:"transaction_time"`
	TransactionType   string              `json:"transaction_type"`
	VaNumbers         []map[string]string `json:"va_numbers"`
	Store             *string             `json:"store"`
	PaymentCode       *string             `json:"payment_code"`
	Permata_va_number *string             `json:"permata_va_number"`
}

func CreateTransaction(orderID string, price int64) (*snap.Response, *midtrans.Error) {
	midtrans.ServerKey = config.CurrentMidtransServerKey()
	midtrans.Environment = midtrans.Sandbox
	if !config.IsInDevelopmentStage() {
		midtrans.Environment = midtrans.Production
	}

	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: price,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	// 3. Request create Snap transaction to Midtrans
	return snap.CreateTransaction(req)
}

type DeterminePaymentSourceResp struct {
	Acquire            string  `json:"acquire"`
	Issuer             *string `json:"issuer"`
	Payment_references *string `json:"payment_references"`
}

func DeterminePaymentSource(midtransNotifResp MidtransNotification) (*DeterminePaymentSourceResp, error) {
	switch midtransNotifResp.PaymentType {
	case "bank_transfer":
		if len(midtransNotifResp.VaNumbers) > 0 {
			return &DeterminePaymentSourceResp{
				Acquire:            midtransNotifResp.VaNumbers[0]["bank"],
				Issuer:             config.String(midtransNotifResp.VaNumbers[0]["bank"]),
				Payment_references: config.String(midtransNotifResp.VaNumbers[0]["va_number"]),
			}, nil
		} else if midtransNotifResp.Permata_va_number != nil {
			return &DeterminePaymentSourceResp{
				Acquire:            "permata",
				Issuer:             config.String("permata"),
				Payment_references: midtransNotifResp.Permata_va_number,
			}, nil
		} else {
			return nil, fmt.Errorf("midtransNotifResponse va number format is not valid")
		}
	case "qris":
		return &DeterminePaymentSourceResp{
			Acquire:            midtransNotifResp.Acquirer,
			Issuer:             midtransNotifResp.Issuer,
			Payment_references: midtransNotifResp.ReferenceID,
		}, nil
	case "echannel":
		return &DeterminePaymentSourceResp{
			Acquire:            *midtransNotifResp.BillKey,
			Issuer:             midtransNotifResp.Issuer,
			Payment_references: midtransNotifResp.ReferenceID,
		}, nil
	case "cstore":
		return &DeterminePaymentSourceResp{
			Acquire:            *midtransNotifResp.Store,
			Issuer:             midtransNotifResp.Store,
			Payment_references: midtransNotifResp.PaymentCode,
		}, nil
	case "dana", "kredivo", "akulaku":
		return &DeterminePaymentSourceResp{
			Acquire:            midtransNotifResp.PaymentType,
			Issuer:             &midtransNotifResp.PaymentType,
			Payment_references: nil,
		}, nil
	default:
		return nil, fmt.Errorf("payment method is not found")
	}
}
