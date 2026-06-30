package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/internal/repository"
	jwtpkg "github.com/topmentor/backend/pkg/jwt"
	smspkg "github.com/topmentor/backend/pkg/sms"
	wxpkg "github.com/topmentor/backend/pkg/wx"
)

type mockUserStore struct {
	byOpenID map[string]*model.User
	byPhone  map[string]*model.User
	nextID   int64
}

func (m *mockUserStore) FindByOpenID(_ context.Context, openID string) (*model.User, error) {
	u, ok := m.byOpenID[openID]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (m *mockUserStore) FindByPhone(_ context.Context, phone string) (*model.User, error) {
	u, ok := m.byPhone[phone]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserStore) FindByID(_ context.Context, id int64) (*model.User, error) {
	for _, u := range m.byPhone {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

func (m *mockUserStore) Create(_ context.Context, openID, phone string) (*model.User, error) {
	m.nextID++
	u := &model.User{ID: m.nextID, OpenID: openID, Phone: phone, EnglishLevel: "beginner", ChildAge: 6}
	m.byOpenID[openID] = u
	m.byPhone[phone] = u
	return u, nil
}

func (m *mockUserStore) UpdateProfile(_ context.Context, id int64, req model.UpdateProfileRequest) (*model.User, error) {
	for _, u := range m.byPhone {
		if u.ID == id {
			u.ChildName = req.ChildName
			u.ChildAge = req.ChildAge
			u.EnglishLevel = req.EnglishLevel
			return u, nil
		}
	}
	return nil, repository.ErrUserNotFound
}

type mockOTPStore struct {
	codes map[string]string
}

func (m *mockOTPStore) Send(_ context.Context, phone string) (string, error) {
	if m.codes == nil {
		m.codes = map[string]string{}
	}
	m.codes[phone] = "123456"
	return "123456", nil
}

func (m *mockOTPStore) Verify(_ context.Context, phone, code string) error {
	if m.codes[phone] != code {
		return smspkg.ErrInvalidCode
	}
	delete(m.codes, phone)
	return nil
}

func newTestAuthService(store *mockUserStore) *AuthService {
	return NewAuthService(store, nil, wxpkg.NewClient("", "", true), jwtpkg.NewManager("secret", 24), &mockOTPStore{}, true)
}

func TestAuthService_SmsLogin_NewUser(t *testing.T) {
	store := &mockUserStore{byOpenID: map[string]*model.User{}, byPhone: map[string]*model.User{}}
	svc := newTestAuthService(store)

	_, err := svc.SendSMS(context.Background(), model.SmsSendRequest{Phone: "13800138000"})
	require.NoError(t, err)

	resp, err := svc.SmsLogin(context.Background(), model.SmsLoginRequest{Phone: "13800138000", Code: "123456"})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "13800138000", resp.User.Phone)
}

func TestAuthService_SmsLogin_InvalidCode(t *testing.T) {
	store := &mockUserStore{byOpenID: map[string]*model.User{}, byPhone: map[string]*model.User{}}
	svc := newTestAuthService(store)

	_, err := svc.SmsLogin(context.Background(), model.SmsLoginRequest{Phone: "13800138000", Code: "000000"})
	assert.ErrorIs(t, err, smspkg.ErrInvalidCode)
}

func TestAuthService_WxLogin_NewUser(t *testing.T) {
	store := &mockUserStore{byOpenID: map[string]*model.User{}, byPhone: map[string]*model.User{}}
	svc := newTestAuthService(store)

	resp, err := svc.WxLogin(context.Background(), model.WxLoginRequest{
		Code:  "user1",
		Phone: "13800138000",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "13800138000", resp.User.Phone)
}

func TestUserService_UpdateProfile(t *testing.T) {
	store := &mockUserStore{
		byOpenID: map[string]*model.User{"oid": {ID: 1, OpenID: "oid", Phone: "13800138000"}},
		byPhone:  map[string]*model.User{"13800138000": {ID: 1, OpenID: "oid", Phone: "13800138000"}},
	}
	svc := NewUserService(store)

	user, err := svc.UpdateProfile(context.Background(), 1, model.UpdateProfileRequest{
		ChildName:    "小明",
		ChildAge:     8,
		EnglishLevel: "intermediate",
	})
	require.NoError(t, err)
	assert.Equal(t, "小明", user.ChildName)
}
