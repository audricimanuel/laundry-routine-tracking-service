package service

import (
	"context"
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/logging"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/repository"
)

type (
	LaundryService interface {
		GetLaundryList(ctx context.Context, queryParam model.LaundryQueryParam, userId string) ([]model.LaundryResponse, error)
		AddLaundry(ctx context.Context, userId string, request model.AddLaundryRequest) error
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

func (l *LaundryServiceImpl) AddLaundry(ctx context.Context, userId string, request model.AddLaundryRequest) error {
	log := logging.WithContext(ctx)

	// validate category id(s)
	var categoryIds []string
	for _, item := range request.Items {
		categoryIds = append(categoryIds, item.CategoryId)
	}

	categoriesData, err := l.laundryRepository.GetCategoryById(ctx, userId, categoryIds...)
	if err != nil || len(categoriesData) != len(request.Items) {
		log.Error("invalid categories detected")
		return errorutils.ErrorBadRequest.CustomMessage("invalid category id")
	}

	return l.laundryRepository.AddLaundryData(ctx, userId, request)
}
