package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/pkg/response"
)

type MentorServicer interface {
	List(ctx context.Context, q model.MentorListQuery) (*model.MentorListResponse, error)
	GetByID(ctx context.Context, id int64) (*model.Mentor, error)
}

type MentorHandler struct {
	svc MentorServicer
}

func NewMentorHandler(svc MentorServicer) *MentorHandler {
	return &MentorHandler{svc: svc}
}

func (h *MentorHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	q := model.MentorListQuery{
		School:   c.Query("school"),
		Gender:   c.Query("gender"),
		Tag:      c.Query("tag"),
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svc.List(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}

	response.OK(c, result)
}

func (h *MentorHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}

	mentor, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	if mentor == nil {
		response.Fail(c, http.StatusNotFound, 40401, "学霸不存在")
		return
	}

	response.OK(c, mentor)
}
