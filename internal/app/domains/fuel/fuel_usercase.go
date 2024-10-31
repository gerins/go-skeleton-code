package fuel

import (
	"context"

	"go-skeleton-code/config"
	"go-skeleton-code/pkg/jwt"
)

type usecase struct {
	fuelRepository Repository
	securityConfig config.Security
}

// NewUsecase returns new usecase.
func NewUsecase(securityConfig config.Security, fuelRepository Repository) *usecase {
	return &usecase{
		fuelRepository: fuelRepository,
		securityConfig: securityConfig,
	}
}

func (u *usecase) List(ctx context.Context, req GetRequest) ([]Fuel, int, error) {
	list, totalItem, err := u.fuelRepository.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	return list, totalItem, nil
}

func (u *usecase) Detail(ctx context.Context, req GetRequest) (Fuel, error) {
	_ = jwt.GetPayloadFromContext(ctx)

	detail, err := u.fuelRepository.Detail(ctx, req)
	if err != nil {
		return Fuel{}, err
	}

	return detail, err
}
