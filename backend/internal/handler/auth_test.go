package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/pkg/response"
)

type stubAuthService struct {
	loginResp *model.LoginResponse
	err       error
}

func (s stubAuthService) SendSMS(_ context.Context, _ model.SmsSendRequest) (*model.SmsSendResponse, error) {
	return &model.SmsSendResponse{ExpiresIn: 300, DebugCode: "123456"}, nil
}

func (s stubAuthService) SmsLogin(_ context.Context, _ model.SmsLoginRequest) (*model.LoginResponse, error) {
	return s.loginResp, s.err
}

func (s stubAuthService) WxLogin(_ context.Context, _ model.WxLoginRequest) (*model.LoginResponse, error) {
	return s.loginResp, s.err
}

func TestAuthHandler_WxLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(stubAuthService{
		loginResp: &model.LoginResponse{Token: "jwt-token", User: model.User{ID: 1, Phone: "13800138000"}},
	})

	body, _ := json.Marshal(model.WxLoginRequest{Code: "abc", Phone: "13800138000"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/wx-login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.WxLogin(c)

	require.Equal(t, http.StatusOK, w.Code)
	var resp response.Body
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 0, resp.Code)
}

func TestAuthHandler_WxLogin_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAuthHandler(stubAuthService{})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/wx-login", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	h.WxLogin(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

type stubMentorService struct {
	listResp *model.MentorListResponse
	getResp  *model.Mentor
}

func (s stubMentorService) List(_ context.Context, _ model.MentorListQuery) (*model.MentorListResponse, error) {
	return s.listResp, nil
}

func (s stubMentorService) GetByID(_ context.Context, id int64) (*model.Mentor, error) {
	return s.getResp, nil
}

func TestMentorHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewMentorHandler(stubMentorService{
		listResp: &model.MentorListResponse{
			List:  []model.Mentor{{ID: 1, RealName: "张同学"}},
			Total: 1,
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/mentors", nil)

	h.List(c)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestMentorHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewMentorHandler(stubMentorService{getResp: nil})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/mentors/99", nil)
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	h.GetByID(c)

	require.Equal(t, http.StatusNotFound, w.Code)
}
