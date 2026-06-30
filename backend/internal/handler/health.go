package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/topmentor/backend/internal/service"
	"github.com/topmentor/backend/pkg/response"
)

type HealthChecker interface {
	Check(ctx context.Context) service.HealthResult
}

type HealthHandler struct {
	svc HealthChecker
}

func NewHealthHandler(svc HealthChecker) *HealthHandler {
	return &HealthHandler{svc: svc}
}

func (h *HealthHandler) Check(c *gin.Context) {
	result := h.svc.Check(c.Request.Context())

	if result.IsHealthy() {
		response.OK(c, result)
		return
	}

	c.JSON(http.StatusServiceUnavailable, response.Body{
		Code: 50001,
		Msg:  "service degraded",
		Data: result,
	})
}
