package fuel

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	response "go-skeleton-code/pkg/response/gin"
)

type httpHandler struct {
	timeout time.Duration
	usecase Usecase
}

func NewHandler(usecase Usecase, timeout time.Duration) interface{ InitRoutes(g *gin.RouterGroup) } {
	return &httpHandler{
		timeout: timeout,
		usecase: usecase,
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

	detail, err := h.usecase.Detail(ctx, requestPayload)
	if err != nil {
		response.Failed(g, err)
		return
	}

	response.Success(g, detail)
}
