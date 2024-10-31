package response

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	serverError "go-skeleton-code/pkg/error"
)

type DefaultResponse struct {
	Code    int `json:"code"`
	Message any `json:"message"`
	Data    any `json:"data"`
	Meta    any `json:"meta,omitempty"`
}

type meta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	MaxPage   int `json:"maxPage"`
	TotalItem int `json:"totalItem"`
}

// Success sends a successful JSON response.
func Success(g *gin.Context, data any) {
	response := DefaultResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	}
	g.JSON(http.StatusOK, response)
}

// SuccessList sends a paginated successful JSON response.
func SuccessList(g *gin.Context, data any, page, limit, totalItem int) {
	response := DefaultResponse{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
		Meta: meta{
			Page:      page,
			Limit:     limit,
			MaxPage:   int(math.Ceil(float64(totalItem) / float64(limit))),
			TotalItem: totalItem,
		},
	}
	g.JSON(http.StatusOK, response)
}

// Failed sends an error response with appropriate HTTP status code.
func Failed(g *gin.Context, err error) {
	var (
		generalError = serverError.ErrGeneralError(nil)
		httpRespCode = generalError.HTTPCode
	)

	// Default response
	response := DefaultResponse{
		Code:    generalError.Code,
		Message: generalError.Message,
		Data:    nil,
	}

	// Default response for data not found in database
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = serverError.ErrDataNotFound(err)
	}

	// Check for wrapped server error
	if serverErr, ok := err.(serverError.ServerError); ok {
		httpRespCode = serverErr.HTTPCode
		response.Code = serverErr.Code
		response.Message = serverErr.Message
	}

	g.JSON(httpRespCode, response)
}
