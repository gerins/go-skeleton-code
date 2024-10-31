package fuel

import (
	"context"

	"gorm.io/gorm"

	gormPkg "go-skeleton-code/pkg/gorm"
	"go-skeleton-code/pkg/log"
)

type repository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewRepository returns new fuel Repository.
func NewRepository(readDB *gorm.DB, writeDB *gorm.DB) *repository {
	return &repository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *repository) List(ctx context.Context, req GetRequest) ([]Fuel, int, error) {
	var (
		totalData int64
		result    = make([]Fuel, 0)
		tableName = Fuel{}.TableName()
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

func (r *repository) Detail(ctx context.Context, req GetRequest) (Fuel, error) {
	var (
		query = r.readDB.WithContext(ctx)
	)

	defer log.Context(ctx).RecordDuration("Detail").Stop()

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	var result Fuel
	if err := query.First(&result, req.ID).Error; err != nil {
		log.Context(ctx).Error(err)
		return Fuel{}, err
	}

	return result, nil
}

func (r *repository) Create(ctx context.Context, payload Fuel) (Fuel, error) {
	defer log.Context(ctx).RecordDuration("Create").Stop()

	if err := r.writeDB.WithContext(ctx).Create(&payload).Error; err != nil {
		log.Context(ctx).Error(err)
		return Fuel{}, err
	}

	return payload, nil
}

func (r *repository) Update(ctx context.Context, payload Fuel) error {
	var (
		query = r.writeDB.WithContext(ctx).Table(Fuel{}.TableName())
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

func (r *repository) Delete(ctx context.Context, id int) error {
	var (
		query = r.writeDB.WithContext(ctx).Where("id", id)
	)

	defer log.Context(ctx).RecordDuration("Delete").Stop()

	exec := query.Delete(&Fuel{})
	if exec.Error != nil {
		log.Context(ctx).Error(exec.Error)
		return exec.Error
	}

	if exec.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
