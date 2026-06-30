package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/topmentor/backend/internal/middleware"
	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/internal/service"
	"github.com/topmentor/backend/pkg/response"
)

type BookingServicer interface {
	ListMentorSlots(ctx context.Context, mentorID int64, fromDate, toDate string) ([]model.MentorSlot, error)
	CreateBooking(ctx context.Context, userID int64, req model.CreateBookingRequest) (*model.BookingResponse, error)
}

type BookingHandler struct {
	svc BookingServicer
}

func NewBookingHandler(svc BookingServicer) *BookingHandler {
	return &BookingHandler{svc: svc}
}

func (h *BookingHandler) ListMentorSlots(c *gin.Context) {
	mentorID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	slots, err := h.svc.ListMentorSlots(c.Request.Context(), mentorID, c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, slots)
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}

	var req model.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.CreateBooking(c.Request.Context(), userID, req)
	if errors.Is(err, repository.ErrSlotNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "时间槽不存在")
		return
	}
	if errors.Is(err, repository.ErrSlotUnavailable) {
		response.Fail(c, http.StatusConflict, 40901, "时间槽已被预约")
		return
	}
	if errors.Is(err, repository.ErrInsufficientLessons) {
		response.Fail(c, http.StatusConflict, 40902, "课时余额不足")
		return
	}
	if errors.Is(err, service.ErrSlotLocked) {
		response.Fail(c, http.StatusConflict, 40901, "时间槽已被占用，请稍后重试")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "预约失败")
		return
	}

	response.OK(c, result)
}

func parseIDParam(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return 0, false
	}
	return id, true
}
