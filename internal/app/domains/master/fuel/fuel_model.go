package fuel

import (
	"context"
	"time"
)

type Fuel struct {
	ID          int        `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement"`
	Type        string     `json:"type" gorm:"column:type;type:varchar;size:255"`
	Description string     `json:"description" gorm:"column:description;type:varchar;size:255"`
	Unit        string     `json:"unit" gorm:"column:unit;type:varchar;size:255"` // Satuan
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime"`
}

func (Fuel) TableName() string {
	return "fuels"
}

type Usecase interface {
	List(ctx context.Context, req GetRequest) ([]Fuel, int, error)
	Detail(ctx context.Context, req GetRequest) (Fuel, error)
}

type Repository interface {
	List(ctx context.Context, req GetRequest) ([]Fuel, int, error)
	Detail(ctx context.Context, req GetRequest) (Fuel, error)
	Create(ctx context.Context, payload Fuel) (Fuel, error)
	Update(ctx context.Context, payload Fuel) error
	Delete(ctx context.Context, id int) error
}
