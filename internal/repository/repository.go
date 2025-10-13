package repository

import (
	"fmt"
	"log"
	"os"
	"time"
	"tsv-golang/pkg/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repositories struct {
	DB      *gorm.DB
	User    UserRepositoryInterface
	Board   BoardRepositoryInterface
	Project ProjectRepositoryInterface
}

func NewRepositories() (*Repositories, error) {
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		dbHost, dbUser, dbPass, dbDatabase, dbPort)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second * 10, // Slow SQL threshold
			LogLevel:                  logger.Silent,    // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,             // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	if env.Env().IsModeDebug() {
		db = db.Debug()
	}
	return &Repositories{
		DB:      db,
		User:    UserRepositoryInit(db),
		Board:   BoardRepositoryInit(db),
		Project: ProjectRepositoryInit(db),
	}, nil
}

func (s *Repositories) Close() error {
	dbInstance, _ := s.DB.DB()
	err := dbInstance.Close()
	return err
}
