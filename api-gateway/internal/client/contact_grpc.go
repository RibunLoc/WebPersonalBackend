package client

import (
	"context"
	"time"

	contactv1 "github.com/RibunLoc/WebPersonalBackend/gen/contact/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type ContactGRPC struct {
	cl contactv1.ContactServiceClient
}

func NewContactGRPC(addr string) (*ContactGRPC, func() error, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 2 * time.Second,
			Backoff: backoff.Config{
				BaseDelay:  200 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   3 * time.Second,
			},
		}),
	)
	if err != nil {
		return nil, nil, err
	}
	return &ContactGRPC{cl: contactv1.NewContactServiceClient(conn)}, conn.Close, nil
}

type ContactSubmitInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Message  string `json:"message"`
	Token    string `json:"turnstile_token,omitempty"`
	CFToken  string `json:"cf_turnstile_token,omitempty"`
	RemoteIP string `json:"-"` // sẽ set từ header
}

type ContactSubmitResult struct {
	Status string `json:"status"`
}

func (c *ContactGRPC) Submit(ctx context.Context, in ContactSubmitInput) (*ContactSubmitResult, error) {
	tok := in.Token
	if tok == "" {
		tok = in.CFToken
	}

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	res, err := c.cl.Submit(ctx, &contactv1.ContactRequest{
		Name:           in.Name,
		Email:          in.Email,
		Message:        in.Message,
		TurnstileToken: tok,
		RemoteIp:       in.RemoteIP,
	})
	if err != nil {
		return nil, err
	}
	return &ContactSubmitResult{Status: res.Status}, nil
}
