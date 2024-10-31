package order

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"go-skeleton-code/config"
	"go-skeleton-code/internal/app/domains/order/model"
	middleware "go-skeleton-code/internal/app/middleware/http/gin"
	response "go-skeleton-code/pkg/response/gin"
)

type httpHandler struct {
	timeout        time.Duration
	orderUsecase   model.Usecase
	securityConfig config.Security
}

func NewHTTPHandler(orderUsecase model.Usecase, timeout time.Duration, securityConfig config.Security) interface {
	InitRoutes(g *gin.RouterGroup)
} {
	return &httpHandler{
		timeout:        timeout,
		orderUsecase:   orderUsecase,
		securityConfig: securityConfig,
	}
}

func (h *httpHandler) InitRoutes(g *gin.RouterGroup) {
	v1 := g.Group("/v1/order")
	v1.Use(middleware.ValidateJwtToken([]byte(h.securityConfig.Jwt.Key)))
	{
		v1.POST("", h.OrderHandler)
	}
}

func (h *httpHandler) OrderHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	var requestPayload model.OrderRequest
	if err := c.Bind(&requestPayload); err != nil {
		response.Failed(c, err)
		return
	}

	orderResult, err := h.orderUsecase.ProcessOrder(ctx, requestPayload)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, orderResult)
}
