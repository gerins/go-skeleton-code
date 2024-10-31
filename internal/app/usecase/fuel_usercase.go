package usecase

import (
	"context"

	"go-skeleton-code/config"
	"go-skeleton-code/internal/app/dto"
	"go-skeleton-code/internal/app/model"
	"go-skeleton-code/pkg/jwt"
)

type usecase struct {
	fuelRepository model.FuelRepository
	securityConfig config.Security
}

// NewFuelUsecase returns new usecase.
func NewFuelUsecase(securityConfig config.Security, fuelRepository model.FuelRepository) *usecase {
	return &usecase{
		fuelRepository: fuelRepository,
		securityConfig: securityConfig,
	}
}

func (u *usecase) List(ctx context.Context, req dto.FuelGetRequest) ([]model.Fuel, int, error) {
	list, totalItem, err := u.fuelRepository.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return list, totalItem, nil
}

func (u *usecase) Detail(ctx context.Context, req dto.FuelGetRequest) (model.Fuel, error) {
	_ = jwt.GetPayloadFromContext(ctx)

	detail, err := u.fuelRepository.Detail(ctx, req)
	if err != nil {
		return model.Fuel{}, err
	}

	return detail, err
}
