package fuel

type GetRequest struct {
	ID   int    `param:"id"`
	Type string `json:"type"`

	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Sort      string `query:"sort"`
	Direction string `validate:"omitempty,oneof=ASC DESC" query:"direction"`
}
