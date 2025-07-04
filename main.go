package main

import (
	"POS-BE/middlewares"
	"POS-BE/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Public routes
	r.POST("/login", services.Login)
	r.POST("/register", services.Register)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middlewares.Authorizer())
	{
		// Categories
		auth.POST("/categories", services.CreateCategories)
		auth.GET("/categories", services.GetCategories)
		auth.PUT("/categories/:id", services.UpdateCategories)
		auth.DELETE("/categories/:id", services.DeleteCategories)

		// Products
		auth.POST("/products", services.CreateProducts)
		auth.GET("/products", services.GetProducts)
		auth.PUT("/products/:id", services.UpdateProducts)
		auth.DELETE("/products/:id", services.DeleteProducts)

		// Transactions
		auth.POST("/transactions/start", services.StartTransaction)
		auth.GET("/transactions", services.GetTransactions)
		auth.GET("/transactions/:id", services.GetTransactionsByID)

		// Transaction Products
		auth.POST("/transactionProducts", services.CreateTransactionProduct)
		auth.GET("/transactionProducts", services.GetTransactionProduct)
		auth.PUT("/transactionProducts/:id", services.UpdateTransactionProduct)
		auth.DELETE("/transactionProducts/:id", services.DeleteTransactionProduct)

		// Payments
		auth.POST("/payments", services.CreatePayment)
	}

	r.Run(":4000")
}
