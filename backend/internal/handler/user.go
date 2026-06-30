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

type UserServicer interface {
	UpdateProfile(ctx context.Context, userID int64, req model.UpdateProfileRequest) (*model.User, error)
}

type UserHandler struct {
	svc UserServicer
}

func NewUserHandler(svc UserServicer) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录或 Token 失效")
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
	if errors.Is(err, repository.ErrUserNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "用户不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "更新失败")
		return
	}

	response.OK(c, user)
}
