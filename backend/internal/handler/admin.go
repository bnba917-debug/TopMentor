package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	"github.com/topmentor/backend/internal/service"
	"github.com/topmentor/backend/pkg/response"
)

type AdminServicer interface {
	Login(ctx context.Context, req model.AdminLoginRequest) (*model.AdminLoginResponse, error)
	ListPendingMentors(ctx context.Context) ([]model.PendingMentorApplication, error)
	ReviewMentor(ctx context.Context, mentorID int64, req model.ReviewMentorRequest) error
	ListCourseware(ctx context.Context) ([]model.Courseware, error)
	CreateCourseware(ctx context.Context, req model.CreateCoursewareRequest) (*model.Courseware, error)
	UpdateCourseware(ctx context.Context, id int64, req model.UpdateCoursewareRequest) (*model.Courseware, error)
	DeleteCourseware(ctx context.Context, id int64) error
	FinanceSummary(ctx context.Context) (*model.FinanceSummary, error)
}

type AdminHandler struct {
	svc AdminServicer
}

func NewAdminHandler(svc AdminServicer) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req model.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	result, err := h.svc.Login(c.Request.Context(), req)
	if errors.Is(err, service.ErrAdminInvalidCredentials) {
		response.Fail(c, http.StatusUnauthorized, 40102, "账号或密码错误")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "登录失败")
		return
	}
	response.OK(c, result)
}

func (h *AdminHandler) ListPendingMentors(c *gin.Context) {
	list, err := h.svc.ListPendingMentors(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *AdminHandler) ReviewMentor(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	var req model.ReviewMentorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	if req.Action == "reject" && req.RejectReason == "" {
		response.Fail(c, http.StatusBadRequest, 40001, "请填写驳回原因")
		return
	}
	err = h.svc.ReviewMentor(c.Request.Context(), id, req)
	if errors.Is(err, service.ErrRejectReasonRequired) {
		response.Fail(c, http.StatusBadRequest, 40001, "请填写驳回原因")
		return
	}
	if errors.Is(err, repository.ErrMentorApplicationNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "申请不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "审核失败")
		return
	}
	response.OK(c, gin.H{"mentor_id": id, "action": req.Action})
}

func (h *AdminHandler) ListCourseware(c *gin.Context) {
	list, err := h.svc.ListCourseware(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *AdminHandler) CreateCourseware(c *gin.Context) {
	var req model.CreateCoursewareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	item, err := h.svc.CreateCourseware(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "创建失败")
		return
	}
	response.OK(c, item)
}

func (h *AdminHandler) UpdateCourseware(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	var req model.UpdateCoursewareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	item, err := h.svc.UpdateCourseware(c.Request.Context(), id, req)
	if errors.Is(err, repository.ErrCoursewareNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "课件不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "更新失败")
		return
	}
	response.OK(c, item)
}

func (h *AdminHandler) DeleteCourseware(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	if err := h.svc.DeleteCourseware(c.Request.Context(), id); errors.Is(err, repository.ErrCoursewareNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "课件不存在")
		return
	} else if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "删除失败")
		return
	}
	response.OK(c, gin.H{"deleted": id})
}

func (h *AdminHandler) FinanceSummary(c *gin.Context) {
	summary, err := h.svc.FinanceSummary(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, summary)
}
