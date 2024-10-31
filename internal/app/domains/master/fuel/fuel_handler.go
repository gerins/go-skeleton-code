package fuel

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	response "go-skeleton-code/pkg/response/gin"
)

type httpHandler struct {
	timeout   time.Duration
	validator *validator.Validate
	usecase   Usecase
}

func NewHandler(timeout time.Duration, validator *validator.Validate, usecase Usecase) interface{ InitRoutes(g *gin.RouterGroup) } {
	return &httpHandler{
		timeout:   timeout,
		validator: validator,
		usecase:   usecase,
	}
}

func (h *httpHandler) InitRoutes(g *gin.RouterGroup) {
	fuel := g.Group("/fuel")
	{
		fuel.GET("", h.GetHandler)
		fuel.GET("/:id", h.DetailHandler)
	}
}

func (h *httpHandler) GetHandler(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), h.timeout)
	defer cancel()

	var requestPayload GetRequest
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

func (h *httpHandler) DetailHandler(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), h.timeout)
	defer cancel()

	var requestPayload GetRequest
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
