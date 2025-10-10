package repository

import "gorm.io/gorm"

type SampleRepository struct {
	db *gorm.DB
}

type SampleRepositoryInterface interface {
}

func SampleRepositoryInit(db *gorm.DB) *SampleRepository {
	return &SampleRepository{
		db: db,
	}
}

var _ SampleRepositoryInterface = &SampleRepository{}
