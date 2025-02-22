package services

import (
	"POS-BE/libraries/config"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = config.InitGormConfig()
}
