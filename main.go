package main

import (
	"POS-BE/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// albums slice to seed record album data.
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//albums RESTAPI
	r.GET("/albums", services.GetAlbums)
	r.POST("/albums", services.PostAlbums)
	r.PUT("/albums/:id", services.PutAlbums)
	r.DELETE("/albums/:id", services.DeleteAlbums)

	//transactions RESTAPI
	r.POST("/start-transactions", services.StartTransaction)
	r.GET("/transactions", services.GetTransactions)

	// r.Run("localhost:4000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(":4000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
