package main

import (
	"POS-BE/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// albums slice to seed record album data.
	r := gin.Default()

	//categories RESTAPI
	r.POST("/categories", services.CreateCategories)
	r.GET("/categories", services.GetCategories)
	r.PUT("/categories/:id", services.UpdateCategories)
	r.DELETE("/categories/:id", services.DeleteCategories)

	//products RESTAPI
	r.POST("/products", services.CreateProducts)
	r.GET("/products", services.GetProducts)
	r.PUT("/products/:id", services.UpdateProducts)
	r.DELETE("/products/:id", services.DeleteProducts)

	//transactions RESTAPI
	r.POST("/transactions/start", services.StartTransaction)
	r.GET("/transactions", services.GetTransactions)
	r.GET("/transactions/:id", services.GetTransactionsByID)

	//payments RESTAPI
	r.POST("/payments", services.CreatePayment)

	r.Run(":4000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
