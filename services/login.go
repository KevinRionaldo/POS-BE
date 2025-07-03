package services

import (
	"POS-BE/libraries/helpers/api/apiResponse"
	"POS-BE/libraries/helpers/utils/authorizer"
	"POS-BE/libraries/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !authorizer.CheckPassword(user.Password, req.Password) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := authorizer.GenerateJWT(user.User_id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": apiResponse.SuccessSingularResponse(map[string]interface{}{
		"token": token,
		"user": models.User{
			User_id: user.User_id,
			Name:    user.Name,
			Email:   user.Email,
			Role:    user.Role,
		}})})
}

func Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Err(err).Msg("Invalid input for registration")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": apiResponse.GeneralErrorResponse(err)})
		return
	}

	// Check if email exists
	var existing models.User
	if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	newUser := models.User{
		User_id:  string(uuid.NewString()),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
		// Role:      "cashier", // or "admin" if you want default role
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
