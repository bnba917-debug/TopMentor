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

type MentorApplyServicer interface {
	GetStatus(ctx context.Context, userID int64) (*model.MentorApplyStatus, error)
	Submit(ctx context.Context, userID int64, req model.SubmitMentorApplyRequest) (*model.MentorApplyStatus, error)
}

type MentorApplyHandler struct {
	svc     MentorApplyServicer
	uploads *upload.Store
}

func NewMentorApplyHandler(svc MentorApplyServicer, uploads *upload.Store) *MentorApplyHandler {
	return &MentorApplyHandler{svc: svc, uploads: uploads}
}

func (h *MentorApplyHandler) GetStatus(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}
	status, err := h.svc.GetStatus(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "查询失败")
		return
	}
	response.OK(c, status)
}

func (h *MentorApplyHandler) Submit(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录")
		return
	}
	var req model.SubmitMentorApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, 40001, "参数校验失败")
		return
	}
	status, err := h.svc.Submit(c.Request.Context(), userID, req)
	if errors.Is(err, repository.ErrAlreadyVerifiedMentor) {
		response.Fail(c, http.StatusConflict, 40906, "您已是认证学霸，请重新登录进入工作台")
		return
	}
	if errors.Is(err, repository.ErrApplicationPending) {
		response.Fail(c, http.StatusConflict, 40907, "申请审核中，请勿重复提交")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 50001, "提交失败")
		return
	}
	response.OK(c, status)
}

func (h *MentorApplyHandler) Upload(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, 40101, "未登录")
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

	uploadKind, ok := applicantUploadKind(kind)
	if !ok {
		response.Fail(c, http.StatusBadRequest, 40001, "无效的上传类型")
		return
	}

	url, err := h.uploads.SaveApplicantFile(userID, uploadKind, file)
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
	response.OK(c, model.UploadResult{URL: url})
}

func applicantUploadKind(kind string) (string, bool) {
	switch kind {
	case "avatar":
		return upload.KindAvatar, true
	case "intro_video":
		return upload.KindIntroVideo, true
	case "id_card":
		return upload.KindIDCard, true
	case "student_card":
		return upload.KindStudentCard, true
	case "english_proof":
		return upload.KindEnglishProof, true
	default:
		return "", false
	}
}
