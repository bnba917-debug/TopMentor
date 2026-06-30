package repository

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/config"
	"github.com/topmentor/backend/internal/model"
	"github.com/topmentor/backend/pkg/database"
)

func TestMentorRepository_UpdateProfile_Integration(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Skip(err)
	}
	if os.Getenv("DB_HOST") == "" && cfg.DBHost == "localhost" {
		// allow local docker default
	}
	pg, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		t.Skip(err)
	}
	defer pg.Close()

	repo := NewMentorRepository(pg.DB())
	ctx := context.Background()

	p, err := repo.UpdateProfile(ctx, 1, model.UpdateMentorProfileRequest{
		RealName:      "张明",
		SchoolName:    "清华大学",
		Major:         "计算机科学",
		Gender:        "male",
		EnglishScore:  "高考英语 148 分",
		Bio:           "集成测试简介",
		Tags:          nil,
		AvatarURL:     "",
		IntroVideoURL: "https://example.com/videos/mentor1.mp4",
	})
	require.NoError(t, err)
	require.Equal(t, "集成测试简介", p.Bio)
}
