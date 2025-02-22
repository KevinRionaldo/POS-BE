package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "github.com/joho/godotenv"
)

var stage = "dev"

func init() {
	if len(os.Args) > 1 {
		stage = os.Args[1]
	}

	envFile := fmt.Sprintf("env/%s.env", stage)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Err(err).Msg("error init env")
	}
	fmt.Printf("running in stage: %s\n", stage)
}

// IsInDevelopmentStage will check the STAGE environment variable is in development stage
func IsInDevelopmentStage() bool {
	return strings.HasSuffix(os.Getenv("STAGE"), "dev")
}

// String converts a string literal to a pointer to that string.
func String(s string) *string {
	return &s
}

// Time converts a time.TIme literal to a pointer to that time.Time.
func Time(t time.Time) *time.Time {
	return &t
}

// CurrentStage return current STAGE environment variable
func CurrentStage() string {
	return os.Getenv("STG")
}

// GetMTMURL returns the MTM_URL environment variable
func GetHeartBeatTimeOut() int {
	valueInt, _ := strconv.Atoi(os.Getenv("HEARTBEAT_TIMEOUT"))
	return valueInt
}

// CurrentMidtransServerKey return midtrans server key for access payment gateway service of midtrans
func CurrentMidtransServerKey() string {
	return os.Getenv("MIDTRANS_SERVER_KEY")
}

// CurrentSchema return current DB_SCHEMA environment variable
func CurrentSchema() string {
	fmt.Println("schema", os.Getenv("STAGE"))
	return os.Getenv("STAGE")
}

// CurrentDatabaseUrl return current DATABASE_URL environment variable
func CurrentDatabaseUrl() string {
	fmt.Println("db url", os.Getenv("DB_URL"))
	return os.Getenv("DB_URL")
}

// GetTableNameOnCurrentSchema return the tableName passing parameter with the current schema
//
//	table := GetTableNameOnCurrentSchema("table1") => it will return [schema].table1
func GetTableNameOnCurrentSchema(tableName string) string {
	targetTableName := fmt.Sprintf("%s.%s", CurrentSchema(), tableName)

	return targetTableName
}

// init connection database for gorm library
func InitGormConfig() *gorm.DB {
	var db *gorm.DB
	dsn := CurrentDatabaseUrl()
	if IsInDevelopmentStage() {
		db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info), TranslateError: true})
	} else {
		db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
	return db
}
