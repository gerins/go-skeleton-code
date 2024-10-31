package user

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	response "go-skeleton-code/pkg/response/gin"
)

type httpHandler struct {
	timeout     time.Duration
	userUsecase Usecase
}

func NewHTTPHandler(userUsecase Usecase, timeout time.Duration) interface{ InitRoutes(g *gin.RouterGroup) } {
	return &httpHandler{
		timeout:     timeout,
		userUsecase: userUsecase,
	}
}

func (h *httpHandler) InitRoutes(g *gin.RouterGroup) {
	v1 := g.Group("/v1/user")
	{
		v1.POST("/login", h.LoginHandler)
		v1.POST("/register", h.RegisterHandler)
	}
}

func (h *httpHandler) LoginHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	var requestPayload LoginRequest
	if err := c.Bind(&requestPayload); err != nil {
		response.Failed(c, err)
		return
	}

	loginResult, err := h.userUsecase.Login(ctx, requestPayload)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, loginResult)
	return
}

func (h *httpHandler) RegisterHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	var requestPayload RegisterRequest
	if err := c.Bind(&requestPayload); err != nil {
		response.Failed(c, err)
		return
	}

	registerResult, err := h.userUsecase.Register(ctx, requestPayload)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, registerResult)
}
