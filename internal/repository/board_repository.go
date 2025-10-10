package repository

import "gorm.io/gorm"

type BoardRepository struct {
	db *gorm.DB
}

type BoardRepositoryInterface interface {
}

func BoardRepositoryInit(db *gorm.DB) *BoardRepository {
	return &BoardRepository{
		db: db,
	}
}

var _ BoardRepositoryInterface = &BoardRepository{}
