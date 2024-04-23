package database

import (
	"github.com/rms_api/internal/config/secrets"
	"github.com/rms_api/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitializeDB(username, password, host, port, dbname string) (*gorm.DB, error) {
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database_test: %v", err)
	}
	if err == nil {
		log.Println("Database is Successfully Initialized", db)
	}
	AutoMigrationErr := db.AutoMigrate(&models.User{}, &models.Token{}, &models.Profile{}, &models.Job{})
	if AutoMigrationErr != nil {
		log.Println("Unable to Migrate Models")
	}
	DB = db
	return DB, nil
}

func OpenDbConnection() (*gorm.DB, error) {
	dbConn, err := InitializeDB(secrets.DbUserName, secrets.DbPassWord, secrets.DbHost, secrets.DbPort, secrets.DbName)
	if err != nil {
		log.Println("Failed to Connect Database, please check DB resource setup.", err)
		return nil, err
	}
	return dbConn, nil
}

func CloseDbConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error getting underlying SQL database:", err)
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		log.Println("Error closing database connection:", err)
	}
	return err
}
