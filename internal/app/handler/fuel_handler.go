package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go-skeleton-code/internal/app/dto"
	"go-skeleton-code/internal/app/model"
	response "go-skeleton-code/pkg/response/gin"
)

type fuelHandler struct {
	validator *validator.Validate
	usecase   model.FuelUsecase
}

func NewFuelHandler(validator *validator.Validate, usecase model.FuelUsecase) interface{ InitRoutes(g *gin.RouterGroup) } {
	return &fuelHandler{
		validator: validator,
		usecase:   usecase,
	}
}

func (h *fuelHandler) InitRoutes(g *gin.RouterGroup) {
	fuel := g.Group("/fuel")
	{
		fuel.GET("", h.GetHandler)
		fuel.GET("/:id", h.DetailHandler)
	}
}

func (h *fuelHandler) GetHandler(g *gin.Context) {
	var (
		ctx            = g.Request.Context()
		requestPayload dto.FuelGetRequest
	)

	if err := g.Bind(&requestPayload); err != nil {
		response.Failed(g, err)
		return
	}

	if err := h.validator.StructCtx(ctx, requestPayload); err != nil {
		response.Failed(g, err)
		return
	}

	list, totalData, err := h.usecase.List(ctx, requestPayload)
	if err != nil {
		response.Failed(g, err)
		return
	}

	response.SuccessList(g, list, requestPayload.Page, requestPayload.Limit, totalData)
}

func (h *fuelHandler) DetailHandler(g *gin.Context) {
	var (
		ctx            = g.Request.Context()
		requestPayload dto.FuelGetRequest
	)

	if err := g.Bind(&requestPayload); err != nil {
		response.Failed(g, err)
		return
	}

	if err := h.validator.StructCtx(ctx, requestPayload); err != nil {
		response.Failed(g, err)
		return
	}

	detail, err := h.usecase.Detail(ctx, requestPayload)
	if err != nil {
		response.Failed(g, err)
		return
	}

	response.Success(g, detail)
}
