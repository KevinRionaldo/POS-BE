package services

import (
	"POS-BE/libraries/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var albums = []models.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func GetAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func PostAlbums(c *gin.Context) {
	var newAlbum models.Album

	log.Info().Any("yyy", "xxx").Msg("test")
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// PutAlbums update an album from JSON received in the request body.
func PutAlbums(c *gin.Context) {
	id := c.Param("id")
	var albumData models.Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&albumData); err != nil {
		return
	}
	albumData.ID = id

	// update album from the slice.
	indexUpdatedData := -1
	for index, item := range albums {
		if item.ID == id {
			indexUpdatedData = index
		}
	}
	if indexUpdatedData == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
	}

	//return response data
	albums[indexUpdatedData] = albumData
	c.IndentedJSON(http.StatusCreated, albumData)
}

// DeleteAlbums delete an album from JSON received in the request body.
func DeleteAlbums(c *gin.Context) {
	id := c.Param("id")
	var albumDeletedData models.Album

	// update album from the slice.
	indexDeletedData := -1
	for index, item := range albums {
		if item.ID == id {
			indexDeletedData = index
		}
	}
	if indexDeletedData == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
	}

	//return response data
	albumDeletedData = albums[indexDeletedData]
	albums = append(albums[:indexDeletedData], albums[indexDeletedData+1:]...)
	c.IndentedJSON(http.StatusCreated, albumDeletedData)
}
