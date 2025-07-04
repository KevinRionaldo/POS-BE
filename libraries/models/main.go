package models

import (
	"POS-BE/libraries/config"
	"time"
)

var schema = config.CurrentSchema()

// Categories model
type Categories struct {
	Categories_id string    `gorm:"primaryKey;not null" json:"categories_id"`
	Name          string    `json:"name"`
	Created_at    time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Created_by    *string   `json:"-"`
	Updated_by    *string   `json:"-"`
}

func (Categories) TableName() string { return schema + ".categories" }

// Product model
type Product struct {
	Product_id    string    `gorm:"primaryKey;not null" json:"product_id"`
	Categories_id *string   `json:"categories_id"`
	Name          string    `json:"name"`
	Stock         int64     `json:"stock"`
	Price         float64   `json:"price"`
	Image         *string   `json:"image"`
	Created_at    time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Created_by    *string   `json:"-"`
	Updated_by    *string   `json:"-"`
}

func (Product) TableName() string { return schema + ".product" }

// Transaction model
type Transaction struct {
	Transaction_id string    `gorm:"primaryKey;not null" json:"transaction_id"`
	User_id        *string   `json:"user_id"`
	Description    string    `json:"description"`
	Created_at     time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Created_by     *string   `json:"-"`
	Updated_by     *string   `json:"-"`
}

func (Transaction) TableName() string { return schema + ".transaction" }

// Transaction_product model
type Transaction_product struct {
	Transaction_product_id string    `gorm:"primaryKey;not null" json:"transaction_product_id"`
	Transaction_id         string    `json:"transaction_id"`
	Product_id             string    `json:"product_id"`
	Quantity               int64     `json:"quantity"`
	Total_price            float64   `json:"total_price"`
	Created_at             time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at             time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Created_by             *string   `json:"-"`
	Updated_by             *string   `json:"-"`
}

func (Transaction_product) TableName() string { return schema + ".transaction_product" }

// Payment model
type Payment struct {
	Payment_id             string     `gorm:"primaryKey;not null" json:"payment_id"`
	Transaction_id         string     `json:"transaction_id"`
	Payment_method         string     `gorm:"not null" json:"payment_method"`
	Acquire                string     `json:"acquire"`
	Issuer                 *string    `json:"issuer"`
	Payment_references     *string    `json:"payment_references"`
	Amount                 float64    `gorm:"type:numeric(15,2);not null" json:"amount"`
	Currency               string     `gorm:"not null;default:'IDR'" json:"currency"`
	Status                 string     `json:"status"`
	Settlement_time        *time.Time `json:"settlement_time"`
	Expiration_time        *time.Time `json:"expiration_time"`
	Created_at             time.Time  `json:"created_at" gorm:"autoCreateTime"`
	Updated_at             time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Created_by             *string    `json:"-"`
	Updated_by             *string    `json:"-"`
	Gateway_transaction_id string     `json:"gateway_transaction_id"`
}

func (Payment) TableName() string { return schema + ".payment" }

type User struct {
	User_id    string    `gorm:"primaryKey;not null" json:"user_id"`
	Name       *string   `json:"name"`
	Email      string    `gorm:"not null" json:"email"`
	Password   string    `gorm:"not null" json:"password"`
	Role       *string   `json:"role"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Created_by *string   `json:"-"`
	Updated_by *string   `json:"-"`
}

func (User) TableName() string { return schema + ".user" }
