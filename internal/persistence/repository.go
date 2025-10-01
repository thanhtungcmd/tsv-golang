package persistence

import (
	"fmt"
	"log"
	"os"
	"time"
	"tsv-golang/pkg/env"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repositories struct {
	DB *gorm.DB
}

func NewRepositories() (*Repositories, error) {
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbDatabase)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second * 10, // Slow SQL threshold
			LogLevel:                  logger.Silent,    // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,             // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return &Repositories{}, nil
	}
	// set debug query where != production env
	if env.Env().IsModeDebug() {
		db = db.Debug()
	}
	return &Repositories{
		DB: db,
	}, nil
}

func (s *Repositories) Close() error {
	dbInstance, _ := s.DB.DB()
	err := dbInstance.Close()
	return err
}
