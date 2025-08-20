package grpcserver

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/RibunLoc/WebPersonalBackend/contact-service/model"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/proto/contactpb"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/repository"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContactGRPC struct {
	contactpb.UnimplementedContactServiceServer
	Repo     *repository.ContactMongo
	Verifier util.Verifier    // có thể là NoopVerifier khi disable
	Emailer  util.EmailSender // có thể nil nếu chưa cấu hình SMTP
}

func Register(s *grpc.Server, repo *repository.ContactMongo, v util.Verifier, emailer util.EmailSender) {
	contactpb.RegisterContactServiceServer(s, &ContactGRPC{
		Repo:     repo,
		Verifier: v,
		Emailer:  emailer,
	})
}

func (h *ContactGRPC) Submit(ctx context.Context, req *contactpb.ContactRequest) (*contactpb.ContactResponse, error) {
	// 1) Validate + trim
	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)
	message := strings.TrimSpace(req.Message)
	if name == "" || email == "" || len(message) < 5 {
		return nil, status.Error(codes.InvalidArgument, "name/email/message too short")
	}

	// 2) Verify Turnstile nếu có Verifier
	if h.Verifier != nil {
		token := strings.TrimSpace(req.TurnstileToken)
		remoteIP := strings.TrimSpace(req.RemoteIp)
		if _, err := h.Verifier.Verify(ctx, token, remoteIP); err != nil {
			return nil, status.Error(codes.InvalidArgument, "turnstile verification failed")
		}
	}

	// 3) Lưu DB
	now := util.CustomTime(time.Now())
	c := &model.Contact{
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: &now,
		// Nếu model có field IP/UserAgent, set thêm ở đây:
		// IP: req.RemoteIp,
	}
	if err := h.Repo.Create(ctx, c); err != nil {
		return nil, status.Errorf(codes.Internal, "db error: %v", err)
	}

	// 4) Gửi email + log kết quả (đừng nuốt lỗi)
	if h.Emailer != nil {
		body := "New contact message (gRPC):\r\n" +
			"Name: " + c.Name + "\r\n" +
			"Email: " + c.Email + "\r\n" +
			"Message:\r\n" + c.Message + "\r\n" +
			"Time: " + c.CreatedAt.String() + "\r\n"

		if err := h.Emailer.Send("[Contact] "+c.Name, body); err != nil {
			log.Printf("[warn] gRPC send mail failed: %v", err)
			// Không fail request — tuỳ policy của bạn
		} else {
			log.Printf("[email] gRPC sent notification to configured recipient")
		}
	} else {
		log.Println("[email] gRPC emailer is nil (disabled/not configured)")
	}

	return &contactpb.ContactResponse{Status: "ok"}, nil
}
