package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

// ConnectDB - open connection to the db
func ConnectDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(os.Getenv("TXN_SERVICE_DB")))
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to TXN_SERVICE_DB")
}

// CloseDB - close db connection
func CloseDB() error {
	var db *gorm.DB

	pdb, err := db.DB()
	if err != nil {
		return err
	}

	pdb.Close()
	fmt.Println("Closing TXN_SERVICE_DB Connection")
	return nil
}
