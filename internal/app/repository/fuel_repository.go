package repository

import (
	"context"

	"gorm.io/gorm"

	"go-skeleton-code/internal/app/dto"
	"go-skeleton-code/internal/app/model"
	gormPkg "go-skeleton-code/pkg/gorm"
	"go-skeleton-code/pkg/log"
)

type fuelRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewFuelRepository returns new model.fuel Repository.
func NewFuelRepository(readDB *gorm.DB, writeDB *gorm.DB) *fuelRepository {
	return &fuelRepository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *fuelRepository) List(ctx context.Context, req dto.FuelGetRequest) ([]model.Fuel, int, error) {
	var (
		totalData int64
		result    = make([]model.Fuel, 0)
		tableName = model.Fuel{}.TableName()
		query     = r.readDB.WithContext(ctx).Table(tableName).Where(tableName + ".deleted_at IS NULL")
	)

	defer log.Context(ctx).RecordDuration("List").Stop()

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// Count Total Item
	query.Count(&totalData)
	query.Scopes(gormPkg.CreatePaginationQuery(req.Page, req.Limit, req.Sort, req.Direction))

	if err := query.Scan(&result).Error; err != nil {
		log.Context(ctx).Error(err)
		return nil, 0, err
	}

	if totalData == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	return result, int(totalData), nil
}

func (r *fuelRepository) Detail(ctx context.Context, req dto.FuelGetRequest) (model.Fuel, error) {
	var (
		query = r.readDB.WithContext(ctx)
	)

	defer log.Context(ctx).RecordDuration("Detail").Stop()

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	var result model.Fuel
	if err := query.First(&result, req.ID).Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Fuel{}, err
	}

	return result, nil
}

func (r *fuelRepository) Create(ctx context.Context, payload model.Fuel) (model.Fuel, error) {
	defer log.Context(ctx).RecordDuration("Create").Stop()

	if err := r.writeDB.WithContext(ctx).Create(&payload).Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Fuel{}, err
	}

	return payload, nil
}

func (r *fuelRepository) Update(ctx context.Context, payload model.Fuel) error {
	var (
		query = r.writeDB.WithContext(ctx).Table(model.Fuel{}.TableName())
	)

	defer log.Context(ctx).RecordDuration("Update").Stop()

	exec := query.Updates(payload)
	if exec.Error != nil {
		log.Context(ctx).Error(exec.Error)
		return exec.Error
	}

	if exec.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *fuelRepository) Delete(ctx context.Context, id int) error {
	var (
		query = r.writeDB.WithContext(ctx).Where("id", id)
	)

	defer log.Context(ctx).RecordDuration("Delete").Stop()

	exec := query.Delete(&model.Fuel{})
	if exec.Error != nil {
		log.Context(ctx).Error(exec.Error)
		return exec.Error
	}

	if exec.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
