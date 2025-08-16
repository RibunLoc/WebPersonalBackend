package application

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/RibunLoc/microservices-learn/auth-service/internal/grpcserver"
	"github.com/RibunLoc/microservices-learn/auth-service/proto/authpb"
	"github.com/RibunLoc/microservices-learn/auth-service/repository"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type App struct {
	router http.Handler
	rdb    *redis.Client // Used for session key storage
	mgdb   *mongo.Database
	config Config
}

func New(ctx context.Context, config Config) (*App, error) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr:     config.RedisAddress,
			Username: config.RedisUsername,
			Password: config.RedisPassword,
		}),
		mgdb:   mongoClient.Database("demo_db"),
		config: config,
	}
	app.loadRoutes()

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	// Running HTTP server
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	go func() {
		if err := a.startGRPCServer(); err != nil {
			ch <- fmt.Errorf("grpc server error: %w", err)
		}
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return httpServer.Shutdown(ctx)
	}

	return nil
}

func (a *App) startGRPCServer() error {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen on port 50051: %w", err)
	}

	grpcServer := grpc.NewServer()

	authpb.RegisterUserServiceServer(grpcServer, &grpcserver.UserGRPCHandler{
		Repo: &repository.RedisMongo{
			Collection: a.mgdb.Collection("users"),
			JwtSecret:  a.config.JwtSecret,
		},
	})
	fmt.Println("gRPC server started on port 50051")
	return grpcServer.Serve(listen)
}
