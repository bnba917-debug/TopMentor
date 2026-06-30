package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/topmentor/backend/internal/middleware"
	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/pkg/response"
)

type RoomServicer interface {
	ListUserOrders(ctx context.Context, userID int64) ([]model.CourseOrderDetail, error)
	GetOrder(ctx context.Context, userID int64, orderID string) (*model.CourseOrderDetail, error)
	Join(ctx context.Context, userID, mentorID int64, orderID string, role string) (*model.JoinRoomResponse, error)
	Heartbeat(ctx context.Context, userID, mentorID int64, orderID, role string) (*model.HeartbeatResponse, error)
	Complete(ctx context.Context, mentorID int64, orderID string) (*model.CompleteRoomResponse, error)
}

type RoomHandler struct {
	svc RoomServicer
}

func NewRoomHandler(svc RoomServicer) *RoomHandler {
	return &RoomHandler{svc: svc}
}

func (h *RoomHandler) ListUserOrders(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}
	list, err := h.svc.ListUserOrders(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *RoomHandler) GetOrder(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}
	order, err := h.svc.GetOrder(c.Request.Context(), userID, c.Param("id"))
	if errors.Is(err, repository.ErrOrderNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "订单不存在")
		return
	}
	if errors.Is(err, repository.ErrOrderForbidden) {
		response.Fail(c, http.StatusForbidden, 40301, "无权限")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, order)
}

func (h *RoomHandler) Join(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}
	mentorID, _ := middleware.GetMentorID(c)

	var req model.JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.Join(c.Request.Context(), userID, mentorID, c.Param("orderId"), req.Role)
	if errors.Is(err, repository.ErrOrderNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "订单不存在")
		return
	}
	if errors.Is(err, repository.ErrOrderForbidden) {
		response.Fail(c, http.StatusForbidden, 40301, "无权限")
		return
	}
	if errors.Is(err, repository.ErrOrderNotJoinable) {
		response.Fail(c, http.StatusConflict, 40903, "当前订单不可进房")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "进房失败")
		return
	}
	response.OK(c, result)
}

func (h *RoomHandler) Heartbeat(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}

	role := c.DefaultQuery("role", "user")
	mentorID, _ := middleware.GetMentorID(c)
	result, err := h.svc.Heartbeat(c.Request.Context(), userID, mentorID, c.Param("orderId"), role)
	if errors.Is(err, repository.ErrOrderNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "订单不存在")
		return
	}
	if errors.Is(err, repository.ErrOrderForbidden) {
		response.Fail(c, http.StatusForbidden, 40301, "无权限")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "心跳失败")
		return
	}
	response.OK(c, result)
}

func (h *RoomHandler) Complete(c *gin.Context) {
	mentorID, ok := middleware.GetMentorID(c)
	if !ok || mentorID <= 0 {
		response.Fail(c, http.StatusForbidden, 40301, "仅学霸可结束课程")
		return
	}

	result, err := h.svc.Complete(c.Request.Context(), mentorID, c.Param("orderId"))
	if errors.Is(err, repository.ErrOrderNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "订单不存在")
		return
	}
	if errors.Is(err, repository.ErrOrderForbidden) {
		response.Fail(c, http.StatusForbidden, 40301, "无权限")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "下课失败")
		return
	}
	response.OK(c, result)
}
