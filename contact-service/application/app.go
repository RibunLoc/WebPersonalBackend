package application

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/RibunLoc/WebPersonalBackend/contact-service/handler"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/internal/grpcserver"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/repository"
	"github.com/RibunLoc/WebPersonalBackend/contact-service/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type App struct {
	cfg         Config
	router      http.Handler
	repo        *repository.ContactMongo
	emailer     util.EmailSender
	mongoClient *mongo.Client
}

func New(ctx context.Context, config Config) (*App, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	// ping để chắc kết nối
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := mongoClient.Database("contact_db")

	// khởi tạo emailer nếu đủ cấu hình
	var emailer util.EmailSender
	if config.SMTPHost != "" && config.SMTPPort != 0 && config.FromEmail != "" && config.NotifyEmail != "" {
		emailer = util.NewSMTPSender(util.SMTPConfig{
			Host:     config.SMTPHost,
			Port:     int(config.SMTPPort),
			Username: config.SMTPUser,
			Password: config.SMTPPassword,
			From:     config.FromEmail,
			To:       []string{config.NotifyEmail},
		})
	}

	if emailer == nil {
		log.Println("[email] DISABLED")
	} else {
		log.Printf("[email] ENABLED host=%s port=%d from=%s to=%s",
			config.SMTPHost, config.SMTPPort, config.FromEmail, config.NotifyEmail)
	}

	app := &App{
		cfg:         config,
		repo:        repository.NewContactRepo(db),
		emailer:     emailer,
		mongoClient: mongoClient,
	}

	app.loadRoutes()
	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.ServerPort),
		Handler:      a.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Println("Starting contact service server")

	errCh := make(chan error, 2)

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("http: %w", err)
		}
	}()

	// gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}
	grpcSrv := grpc.NewServer()
	verifier := util.NewVerifier(a.cfg.TurnstileSecret, a.cfg.TurnstileDisable)
	grpcserver.Register(grpcSrv, a.repo, verifier, a.emailer)

	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			errCh <- fmt.Errorf("grpc: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		_ = httpSrv.Shutdown(ctx)
		grpcSrv.GracefulStop()
		if a.mongoClient != nil {
			_ = a.mongoClient.Disconnect(ctx)
		}
		return nil
	}
}

func (a *App) buildHandlers() *handler.ContactHandler {
	return &handler.ContactHandler{
		Repo:     a.repo,
		Verifier: util.NewVerifier(a.cfg.TurnstileSecret, a.cfg.TurnstileDisable),
		Emailer:  a.emailer,
	}
}
