package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/topmentor/backend/internal/model"
	smspkg "github.com/topmentor/backend/pkg/sms"
	"github.com/topmentor/backend/pkg/response"
)

type AuthServicer interface {
	SendSMS(ctx context.Context, req model.SmsSendRequest) (*model.SmsSendResponse, error)
	SmsLogin(ctx context.Context, req model.SmsLoginRequest) (*model.LoginResponse, error)
	WxLogin(ctx context.Context, req model.WxLoginRequest) (*model.LoginResponse, error)
}

type AuthHandler struct {
	svc AuthServicer
}

func NewAuthHandler(svc AuthServicer) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) SendSMS(c *gin.Context) {
	var req model.SmsSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.SendSMS(c.Request.Context(), req)
	if errors.Is(err, smspkg.ErrSendTooFrequent) {
		response.Fail(c, http.StatusTooManyRequests, 42901, "发送过于频繁，请稍后再试")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "发送失败")
		return
	}

	response.OK(c, result)
}

func (h *AuthHandler) SmsLogin(c *gin.Context) {
	var req model.SmsLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.SmsLogin(c.Request.Context(), req)
	if errors.Is(err, smspkg.ErrInvalidCode) {
		response.Fail(c, http.StatusBadRequest, 40002, "验证码错误或已过期")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "登录失败")
		return
	}

	response.OK(c, result)
}

func (h *AuthHandler) WxLogin(c *gin.Context) {
	var req model.WxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.WxLogin(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "登录失败")
		return
	}

	response.OK(c, result)
}
