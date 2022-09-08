package api

import (
	"context"
	"time"

	"github.com/Dsmit05/party-day-bot/internal/models"
)

type BotI interface {
	SendMsgAllUser(ctx context.Context, msg string) error
	SendPhoto(ctx context.Context, user models.User, url string) error
	SendMsg(ctx context.Context, user models.User, msg string) error
}

type ConfigI interface {
	GetApiGRPCServerAddr() string
	GetApiGRPCServerTimeout() time.Duration
	GetApiHTTPServerAddr() string
}
