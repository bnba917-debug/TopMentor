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
	"github.com/topmentor/backend/pkg/upload"
)

type MentorPortalServicer interface {
	ListOrders(ctx context.Context, mentorID int64) ([]model.CourseOrderDetail, error)
	ListSlots(ctx context.Context, mentorID int64, from, to string) ([]model.MentorSlot, error)
	SetSlots(ctx context.Context, mentorID int64, req model.SetSlotsRequest) error
	SubmitReport(ctx context.Context, mentorID int64, req model.SubmitReportRequest) (*model.GrowthReport, error)
	GetReportForUser(ctx context.Context, userID int64, orderID string) (*model.GrowthReport, error)
	Wallet(ctx context.Context, mentorID int64) (*model.WalletSummary, error)
	Withdraw(ctx context.Context, mentorID int64, req model.WithdrawRequest) (*model.WithdrawResponse, error)
	GetProfile(ctx context.Context, mentorID int64) (*model.MentorPortalProfile, error)
	UpdateProfile(ctx context.Context, mentorID int64, req model.UpdateMentorProfileRequest) (*model.MentorPortalProfile, error)
	ApplyUploadedMedia(ctx context.Context, mentorID int64, kind, url string) error
}

type MentorPortalHandler struct {
	svc     MentorPortalServicer
	uploads *upload.Store
}

func NewMentorPortalHandler(svc MentorPortalServicer, uploads *upload.Store) *MentorPortalHandler {
	return &MentorPortalHandler{svc: svc, uploads: uploads}
}

func mentorIDFromContext(c *gin.Context) (int64, bool) {
	return middleware.GetMentorID(c)
}

func (h *MentorPortalHandler) requireMentor(c *gin.Context) (int64, bool) {
	id, ok := mentorIDFromContext(c)
	if !ok || id <= 0 {
		response.Fail(c, http.StatusForbidden, 40301, "请使用学霸账号登录")
		return 0, false
	}
	return id, true
}

func (h *MentorPortalHandler) ListOrders(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	list, err := h.svc.ListOrders(c.Request.Context(), mentorID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, list)
}

func (h *MentorPortalHandler) ListSlots(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	slots, err := h.svc.ListSlots(c.Request.Context(), mentorID, c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, slots)
}

func (h *MentorPortalHandler) SetSlots(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	var req model.SetSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	if err := h.svc.SetSlots(c.Request.Context(), mentorID, req); err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "设置失败")
		return
	}
	response.OK(c, gin.H{"updated": len(req.Slots)})
}

func (h *MentorPortalHandler) SubmitReport(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	var req model.SubmitReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	report, err := h.svc.SubmitReport(c.Request.Context(), mentorID, req)
	if errors.Is(err, repository.ErrReportExists) {
		response.Fail(c, http.StatusConflict, 40904, "报告已提交")
		return
	}
	if errors.Is(err, repository.ErrOrderNotCompleted) {
		response.Fail(c, http.StatusBadRequest, 40004, "课程尚未完成")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "提交失败")
		return
	}
	response.OK(c, report)
}

func (h *MentorPortalHandler) Wallet(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	w, err := h.svc.Wallet(c.Request.Context(), mentorID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, w)
}

func (h *MentorPortalHandler) Withdraw(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	var req model.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	result, err := h.svc.Withdraw(c.Request.Context(), mentorID, req)
	if errors.Is(err, repository.ErrInsufficientBalance) {
		response.Fail(c, http.StatusConflict, 40905, "余额不足")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "提现失败")
		return
	}
	response.OK(c, result)
}

func (h *MentorPortalHandler) GetReport(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}
	report, err := h.svc.GetReportForUser(c.Request.Context(), userID, c.Param("orderId"))
	if errors.Is(err, repository.ErrReportNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "报告不存在")
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
	response.OK(c, report)
}

func (h *MentorPortalHandler) GetProfile(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	profile, err := h.svc.GetProfile(c.Request.Context(), mentorID)
	if errors.Is(err, repository.ErrMentorNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "学霸资料不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, profile)
}

func (h *MentorPortalHandler) UpdateProfile(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	var req model.UpdateMentorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败："+err.Error())
		return
	}
	profile, err := h.svc.UpdateProfile(c.Request.Context(), mentorID, req)
	if errors.Is(err, repository.ErrMentorNotFound) {
		response.Fail(c, http.StatusNotFound, 40401, "学霸资料不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "保存失败")
		return
	}
	response.OK(c, profile)
}

func (h *MentorPortalHandler) Upload(c *gin.Context) {
	mentorID, ok := h.requireMentor(c)
	if !ok {
		return
	}
	if h.uploads == nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "上传服务未配置")
		return
	}

	kind := c.PostForm("kind")
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "请选择文件")
		return
	}

	uploadKind := kind
	switch kind {
	case "avatar":
		uploadKind = upload.KindAvatar
	case "intro_video":
		uploadKind = upload.KindIntroVideo
	default:
		response.Fail(c, http.StatusBadRequest, 40001, "无效的上传类型")
		return
	}

	url, err := h.uploads.SaveMentorFile(mentorID, uploadKind, file)
	if errors.Is(err, upload.ErrFileTooLarge) {
		response.Fail(c, http.StatusBadRequest, 40002, "文件过大")
		return
	}
	if errors.Is(err, upload.ErrUnsupportedExt) {
		response.Fail(c, http.StatusBadRequest, 40003, "不支持的文件格式")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "上传失败")
		return
	}

	if err := h.svc.ApplyUploadedMedia(c.Request.Context(), mentorID, kind, url); err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "保存文件地址失败")
		return
	}
	response.OK(c, model.UploadResult{URL: url})
}
