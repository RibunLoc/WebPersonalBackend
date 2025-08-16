package grpcserver

import (
	"context"
	"errors"

	"github.com/RibunLoc/microservices-learn/auth-service/proto/authpb"
	repository "github.com/RibunLoc/microservices-learn/auth-service/repository"

	"github.com/RibunLoc/microservices-learn/auth-service/util"
)

type UserGRPCHandler struct {
	authpb.UnimplementedUserServiceServer
	Repo *repository.RedisMongo
}

// Đăng nhập người dùng
func (h *UserGRPCHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user, err := h.Repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	token, err := util.GenerateJWT(user.ID.Hex(), h.Repo.JwtSecret)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Token: token,
		User: &authpb.UserResponse{
			Id:       user.ID.Hex(),
			Email:    user.Email,
			Fullname: user.Fullname,
			Role:     user.Role,
		},
	}, nil
}
