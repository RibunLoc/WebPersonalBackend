package client

import (
	"context"
	"time"

	"github.com/RibunLoc/WebPersonalBackend/auth-service/proto/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthGRPC struct {
	cl authpb.UserServiceClient
}

func NewAuthGRPC(addr string) (*AuthGRPC, func() error, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	)
	if err != nil {
		return nil, nil, err
	}
	return &AuthGRPC{cl: authpb.NewUserServiceClient(conn)}, conn.Close, nil
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Role     string `json:"role"`
	} `json:"user"`
}

func (a *AuthGRPC) Login(ctx context.Context, in LoginInput) (*LoginResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := a.cl.Login(ctx, &authpb.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}

	var out LoginResult
	out.Token = res.Token
	out.User.ID = res.User.Id
	out.User.Email = res.User.Email
	out.User.Fullname = res.User.Fullname
	out.User.Role = res.User.Role
	return &out, nil
}
