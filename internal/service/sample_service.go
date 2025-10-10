package service

import "tsv-golang/internal/repository"

type SampleService struct {
	repo *repository.Repositories
}

type SampleServiceInterface interface {
}

func SampleServiceInit(repo *repository.Repositories) *SampleService {
	return &SampleService{
		repo: repo,
	}
}

var _ SampleServiceInterface = &SampleService{}
