package main

import (
	"POS-BE/libraries/config"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.InitGormConfig()
}

// InitDB initializes the database and runs migrations with foreign keys
func main() {
	schema := config.CurrentSchema()

	//create schema if don't exists
	if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", config.CurrentSchema())).Error; err != nil {
		log.Fatal(err.Error())
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s."user" (
			user_id VARCHAR NOT NULL PRIMARY KEY,
			name VARCHAR NULL,
			email VARCHAR NOT NULL,
			password VARCHAR NOT NULL,
			"role" VARCHAR NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			created_by VARCHAR NULL,
			updated_by VARCHAR NULL
		);

		CREATE TABLE IF NOT EXISTS %[1]s.categories (
			categories_id VARCHAR NOT NULL PRIMARY KEY,
			name VARCHAR NOT NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			created_by VARCHAR NULL,
			updated_by VARCHAR NULL
		);

		CREATE TABLE IF NOT EXISTS %[1]s.product (
			product_id VARCHAR NOT NULL PRIMARY KEY,
			categories_id VARCHAR NOT NULL,
			name VARCHAR NOT NULL,
			stock INT8 NOT NULL,
			price FLOAT8 NOT NULL,
			image VARCHAR NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			created_by VARCHAR NULL,
			updated_by VARCHAR NULL,
			CONSTRAINT product_categories_fk FOREIGN KEY (categories_id) REFERENCES %[1]s.categories(categories_id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS %[1]s.transaction (
			transaction_id VARCHAR NOT NULL PRIMARY KEY,
			user_id VARCHAR NULL,
			description STRING NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			created_by VARCHAR NULL,
			updated_by VARCHAR NULL,
			CONSTRAINT transaction_user_fk FOREIGN KEY (user_id) REFERENCES %[1]s."user"(user_id) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS %[1]s.payment (
			payment_id VARCHAR NOT NULL PRIMARY KEY,
			transaction_id VARCHAR NULL,
			payment_method VARCHAR NOT NULL,
			acquire VARCHAR NOT NULL,
			issuer VARCHAR NULL,
			payment_references VARCHAR NULL,
			amount FLOAT8 NOT NULL,
			currency VARCHAR NOT NULL DEFAULT 'IDR',
			status VARCHAR NOT NULL,
			settlement_time TIMESTAMP NULL,
			expiration_time TIMESTAMP NULL,
			gateway_transaction_id VARCHAR NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			created_by VARCHAR NULL,
			updated_by VARCHAR NULL,
			CONSTRAINT payment_transaction_fk FOREIGN KEY (transaction_id) REFERENCES %[1]s.transaction(transaction_id)
		);

		CREATE TABLE IF NOT EXISTS %[1]s.transaction_product (
			transaction_product_id VARCHAR NOT NULL PRIMARY KEY,
			product_id VARCHAR NOT NULL,
			transaction_id VARCHAR NOT NULL,
			quantity INT8 NOT NULL,
			total_price FLOAT8 NOT NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			created_by VARCHAR NULL,
			updated_by VARCHAR NULL,
			CONSTRAINT transaction_product_product_fk FOREIGN KEY (product_id) REFERENCES %[1]s.product(product_id) ON DELETE CASCADE,
			CONSTRAINT transaction_product_transaction_fk FOREIGN KEY (transaction_id) REFERENCES %[1]s.transaction(transaction_id) ON DELETE CASCADE,
			UNIQUE INDEX transaction_product_unique (product_id ASC, transaction_id ASC
		);
		`, schema)

		if err := tx.Exec(sql).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	} else {
		log.Println("Database initialization completed successfully.")
	}
}
