package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/service"
	"github.com/topmentor/backend/pkg/response"
)

type stubHealthChecker struct {
	result service.HealthResult
}

func (s stubHealthChecker) Check(_ context.Context) service.HealthResult {
	return s.result
}

func TestHealthHandler_AllUp_Returns200(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewHealthHandler(stubHealthChecker{
		result: service.HealthResult{
			Status: "ok",
			DB:     service.StatusUp,
			Redis:  service.StatusUp,
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)

	h.Check(c)

	require.Equal(t, http.StatusOK, w.Code)

	var body response.Body
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 0, body.Code)
	assert.Equal(t, "ok", body.Msg)
}

func TestHealthHandler_Degraded_Returns503(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewHealthHandler(stubHealthChecker{
		result: service.HealthResult{
			Status: "degraded",
			DB:     service.StatusUp,
			Redis:  service.StatusDown,
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)

	h.Check(c)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)

	var body response.Body
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, 50001, body.Code)
	assert.Equal(t, "service degraded", body.Msg)
}
