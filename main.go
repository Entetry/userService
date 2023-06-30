// Package main
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Entetry/userService/internal/config"
	"github.com/Entetry/userService/internal/handler"
	"github.com/Entetry/userService/internal/repository"
	"github.com/Entetry/userService/internal/service"
	"github.com/Entetry/userService/protocol/userService"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	db, err := pgxpool.Connect(ctx, cfg.ConnectionString)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s\n", err) //nolint:errcheck,gocritic
	}
	defer db.Close()
	userRepository := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepository)
	userHandler := handler.NewUser(userSvc)
	grpcServer := grpc.NewServer()
	userService.RegisterUserServiceServer(grpcServer, userHandler)
	go func() {
		<-sigChan
		cancel()
		grpcServer.GracefulStop()
		if err != nil {
			log.Errorf("can't stop server gracefully %v", err)
		}
	}()
	log.Info("grpc Server started on ", cfg.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
