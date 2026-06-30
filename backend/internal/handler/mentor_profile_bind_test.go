package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
)

func TestUpdateMentorProfileRequest_Binding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := map[string]interface{}{
		"real_name":       "张明",
		"school_name":     "清华大学",
		"major":           "计算机科学",
		"gender":          "male",
		"english_score":   "高考英语 148 分",
		"bio":             "简介",
		"tags":            []string{"阳光幽默", "善于引导"},
		"avatar_url":      "",
		"intro_video_url": "https://example.com/v.mp4",
	}
	raw, _ := json.Marshal(body)

	var req model.UpdateMentorProfileRequest
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPut, "/mentor/profile", bytes.NewReader(raw))
	c.Request.Header.Set("Content-Type", "application/json")
	require.NoError(t, c.ShouldBindJSON(&req))
	assert.Equal(t, "张明", req.RealName)
	assert.Len(t, req.Tags, 2)
}

func TestUpdateMentorProfileRequest_EmptyTags(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := `{"real_name":"a","school_name":"b","major":"c","gender":"unknown","tags":[]}`
	var req model.UpdateMentorProfileRequest
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	err := c.ShouldBindJSON(&req)
	assert.NoError(t, err)
}
