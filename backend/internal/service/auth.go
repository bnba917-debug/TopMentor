package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	jwtpkg "github.com/topmentor/backend/pkg/jwt"
	wxpkg "github.com/topmentor/backend/pkg/wx"
)

type UserStore interface {
	FindByOpenID(ctx context.Context, openID string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, openID, phone string) (*model.User, error)
	UpdateProfile(ctx context.Context, id int64, req model.UpdateProfileRequest) (*model.User, error)
}

type OTPStore interface {
	Send(ctx context.Context, phone string) (string, error)
	Verify(ctx context.Context, phone, code string) error
}

type MentorFinder interface {
	FindByPhone(ctx context.Context, phone string) (*model.MentorProfile, error)
}

type AuthService struct {
	users      UserStore
	mentors    MentorFinder
	wx         *wxpkg.Client
	jwt        *jwtpkg.Manager
	sms        OTPStore
	smsMock    bool
}

func NewAuthService(users UserStore, mentors MentorFinder, wx *wxpkg.Client, jwt *jwtpkg.Manager, sms OTPStore, smsMock bool) *AuthService {
	return &AuthService{users: users, mentors: mentors, wx: wx, jwt: jwt, sms: sms, smsMock: smsMock}
}

func (s *AuthService) SendSMS(ctx context.Context, req model.SmsSendRequest) (*model.SmsSendResponse, error) {
	code, err := s.sms.Send(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	resp := &model.SmsSendResponse{ExpiresIn: 300}
	if s.smsMock {
		resp.DebugCode = code
	}
	return resp, nil
}

func (s *AuthService) SmsLogin(ctx context.Context, req model.SmsLoginRequest) (*model.LoginResponse, error) {
	if err := s.sms.Verify(ctx, req.Phone, req.Code); err != nil {
		return nil, err
	}

	user, err := s.users.FindByPhone(ctx, req.Phone)
	if errors.Is(err, repository.ErrUserNotFound) {
		openID := fmt.Sprintf("sms_%s", req.Phone)
		user, err = s.users.Create(ctx, openID, req.Phone)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return s.issueToken(ctx, user)
}

func (s *AuthService) WxLogin(ctx context.Context, req model.WxLoginRequest) (*model.LoginResponse, error) {
	session, err := s.wx.Code2Session(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByOpenID(ctx, session.OpenID)
	if errors.Is(err, sql.ErrNoRows) {
		user, err = s.users.Create(ctx, session.OpenID, req.Phone)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return s.issueToken(ctx, user)
}

func (s *AuthService) issueToken(ctx context.Context, user *model.User) (*model.LoginResponse, error) {
	var mentorID int64
	var mentor *model.MentorProfile
	if s.mentors != nil {
		if m, err := s.mentors.FindByPhone(ctx, user.Phone); err == nil && m != nil {
			mentor = m
			mentorID = m.ID
		}
	}

	token, err := s.jwt.Issue(user.ID, user.OpenID, mentorID)
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{Token: token, User: *user, Mentor: mentor}, nil
}

type UserService struct {
	users UserStore
}

func NewUserService(users UserStore) *UserService {
	return &UserService{users: users}
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req model.UpdateProfileRequest) (*model.User, error) {
	return s.users.UpdateProfile(ctx, userID, req)
}
