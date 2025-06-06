package service

import (
	"context"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/repository"
)

type (
	LaundryService interface {
		GetLaundryList(ctx context.Context, queryParam model.LaundryQueryParam, userId string) ([]model.LaundryResponse, error)
	}

	LaundryServiceImpl struct {
		laundryRepository repository.LaundryRepository
	}
)

func NewLaundryService(laundryRepo repository.LaundryRepository) LaundryService {
	return &LaundryServiceImpl{
		laundryRepository: laundryRepo,
	}
}

func (l *LaundryServiceImpl) GetLaundryList(ctx context.Context, queryParam model.LaundryQueryParam, userId string) ([]model.LaundryResponse, error) {
	return l.laundryRepository.GetLaundryList(ctx, queryParam, userId)
}
