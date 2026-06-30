package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/topmentor/backend/internal/config"
	"github.com/topmentor/backend/pkg/database"
)

func TestMentorRepository_FindOwnedByID_NullTags(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)
	pg, err := database.NewPostgres(cfg.DSN())
	require.NoError(t, err)
	defer pg.Close()

	repo := NewMentorRepository(pg.DB())
	p, err := repo.FindOwnedByID(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, p)
}
