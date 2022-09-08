package bot

import (
	"context"

	"github.com/Dsmit05/party-day-bot/internal/models"
)

type cacheI interface {
	CheckAccess(ctx context.Context, userID int64) bool
	GetAdmins(ctx context.Context) []int64
	UpdateRole(ctx context.Context, userID int64, role string) error
	GetUsers(ctx context.Context) map[int64]string
	AddUser(ctx context.Context, user models.User) error
	CheckUser(ctx context.Context, userID int64) bool
}

type dbI interface {
	CreateMessage(ctx context.Context, userTgID int64, text string) error
	CreateFile(ctx context.Context, file models.File) error
}
