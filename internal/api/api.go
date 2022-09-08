package api

import (
	"context"

	"github.com/Dsmit05/party-day-bot/internal/logger"
	"github.com/Dsmit05/party-day-bot/internal/models"
	pb "github.com/Dsmit05/party-day-bot/pkg/api"
)

type UserServerAPI struct {
	bot BotI
	pb.UnimplementedUserServer
}

func NewUserServerAPI(bot BotI) pb.UserServer {
	return &UserServerAPI{
		bot: bot,
	}
}

func (a *UserServerAPI) SendMsgAllUser(ctx context.Context, in *pb.SendMsgAllUserRequest) (*pb.SendMsgAllUserResponse, error) {
	if err := a.bot.SendMsgAllUser(ctx, in.GetMsg()); err != nil {
		logger.Error(" a.bot.SendMsgAllUser()", err)
		return nil, err
	}

	return &pb.SendMsgAllUserResponse{}, nil
}

func (a *UserServerAPI) SendPhoto(ctx context.Context, in *pb.SendPhotoRequest) (*pb.SendPhotoResponse, error) {
	user := models.User{
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		UserName:  in.GetUserName(),
	}

	if err := a.bot.SendPhoto(ctx, user, in.GetUrl()); err != nil {
		logger.Error(" a.bot.SendFoto()", err)
		return nil, err
	}

	return &pb.SendPhotoResponse{}, nil
}

func (a *UserServerAPI) SendMsg(ctx context.Context, in *pb.SendMsgRequest) (*pb.SendMsgResponse, error) {
	user := models.User{
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		UserName:  in.GetUserName(),
	}

	if err := a.bot.SendMsg(ctx, user, in.GetMsg()); err != nil {
		logger.Error(" a.bot.SendMsg()", err)
		return nil, err
	}

	return &pb.SendMsgResponse{}, nil
}
