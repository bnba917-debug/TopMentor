package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/topmentor/backend/internal/middleware"
	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/pkg/payment"
	"github.com/topmentor/backend/pkg/response"
)

type RechargeServicer interface {
	ListPackages(ctx context.Context) ([]model.LessonPackage, error)
	PaymentChannels() model.PaymentChannelsResponse
	CreateRecharge(ctx context.Context, userID int64, req model.RechargeRequest, clientIP string) (*model.RechargeResponse, error)
	GetLessonBalance(ctx context.Context, userID int64) (*model.LessonBalanceResponse, error)
	GetRechargeOrder(ctx context.Context, userID int64, orderID string) (*model.RechargeOrder, error)
}

type RechargeHandler struct {
	svc RechargeServicer
}

func NewRechargeHandler(svc RechargeServicer) *RechargeHandler {
	return &RechargeHandler{svc: svc}
}

func (h *RechargeHandler) ListPackages(c *gin.Context) {
	list, err := h.svc.ListPackages(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *RechargeHandler) PaymentChannels(c *gin.Context) {
	response.OK(c, h.svc.PaymentChannels())
}

func (h *RechargeHandler) CreateRecharge(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}

	var req model.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.CreateRecharge(c.Request.Context(), userID, req, c.ClientIP())
	if errors.Is(err, repository.ErrPackageNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "课时包不存在")
		return
	}
	if errors.Is(err, payment.ErrNotConfigured) {
		response.Fail(c, http.StatusBadRequest, 40003, "该支付方式尚未配置，请填写商户号或改用 mock")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "充值失败")
		return
	}

	response.OK(c, result)
}

func (h *RechargeHandler) GetLessonBalance(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}

	balance, err := h.svc.GetLessonBalance(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, balance)
}

func (h *RechargeHandler) GetRechargeOrder(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}

	order, err := h.svc.GetRechargeOrder(c.Request.Context(), userID, c.Param("id"))
	if errors.Is(err, repository.ErrRechargeNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "订单不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, order)
}

// NotifyWechat is a placeholder for WeChat payment async notification (live mode).
func (h *RechargeHandler) NotifyWechat(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, 50101, "微信支付回调待接入")
}
