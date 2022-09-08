package api

import (
	"context"
	"testing"

	mock_api "github.com/Dsmit05/party-day-bot/internal/api/mocks"
	"github.com/Dsmit05/party-day-bot/internal/logger"
	"github.com/Dsmit05/party-day-bot/internal/models"
	pb "github.com/Dsmit05/party-day-bot/pkg/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUserServerApi_SendMsg(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	bot := mock_api.NewMockBotI(ctl)
	ctx := context.Background()
	_ = logger.InitLogger(true, "")

	inData := &pb.SendMsgRequest{
		Msg:       "hi",
		FirstName: "FirstName",
		LastName:  "LastName",
		UserName:  "UserName",
	}

	user := models.User{
		FirstName: "FirstName",
		LastName:  "LastName",
		UserName:  "UserName",
	}

	bot.EXPECT().SendMsg(ctx, user, "hi").Return(nil)

	serverApi := NewUserServerAPI(bot)
	_, err := serverApi.SendMsg(ctx, inData)

	require.NoError(t, err)
}

func TestUserServerApi_SendMsgAllUser(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	bot := mock_api.NewMockBotI(ctl)
	ctx := context.Background()
	_ = logger.InitLogger(true, "")

	inData := &pb.SendMsgAllUserRequest{
		Msg: "hi",
	}

	bot.EXPECT().SendMsgAllUser(ctx, "hi").Return(nil)

	serverApi := NewUserServerAPI(bot)
	_, err := serverApi.SendMsgAllUser(ctx, inData)

	require.NoError(t, err)
}

func TestUserServerApi_SendPhoto(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	bot := mock_api.NewMockBotI(ctl)
	ctx := context.Background()
	_ = logger.InitLogger(true, "")

	inData := &pb.SendPhotoRequest{
		Url:       "test.url/m.img",
		FirstName: "FirstName",
		LastName:  "LastName",
		UserName:  "UserName",
	}

	user := models.User{
		FirstName: "FirstName",
		LastName:  "LastName",
		UserName:  "UserName",
	}

	bot.EXPECT().SendPhoto(ctx, user, "test.url/m.img").Return(nil)

	serverAPI := NewUserServerAPI(bot)
	_, err := serverAPI.SendPhoto(ctx, inData)

	require.NoError(t, err)
}
