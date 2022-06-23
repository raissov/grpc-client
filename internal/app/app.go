package app

import (
	"client/config"
	delivery "client/internal/delivery/http/v1"
	"client/internal/repository"
	"client/internal/service"
	"client/pb"
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	var logger *zap.Logger
	var err error
	appMode := os.Getenv("APP_MODE")
	appMode = "dev"
	if appMode == "prod" {
		logger, err = zap.NewProduction()
		if err != nil {
			log.Println("Error while creating logger: ", err)
			return err
		}
	} else if appMode == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Println("Error while creating logger: ", err)
			return err
		}
	} else {
		log.Println("Error while creating logger: logger mode undefined")
		return fmt.Errorf("error while defining logger: app_mode is invalid %s", appMode)
	}
	//nolint
	defer logger.Sync()
	sugar := logger.Sugar()
	cfg, err := config.New()
	if err != nil {
		sugar.Errorf("errof while defining config: %v", err)
		return err
	}

	sAddress := fmt.Sprintf("%s:%d", cfg.GRPCConfig.Host, cfg.GRPCConfig.Port)

	conn, err := grpc.Dial(sAddress, grpc.WithInsecure())
	if err != nil {
		sugar.Errorf("error while dialing grpc: %v", err)
	}
	defer conn.Close()
	c := pb.NewMessageServiceClient(conn)

	repos := repository.NewRepository(sugar)
	if err != nil {
		sugar.Errorf("error while creating repository: %v", err)
		return err
	}
	services := service.NewService(repos, cfg, sugar)
	handlers := delivery.NewHandler(services, sugar, cfg, c)
	logger.Info(cfg.App.AppPort)
	srv := http.Server{
		Addr:    ":" + cfg.App.AppPort,
		Handler: handlers.InitRoutes(),
	}

	errChan := make(chan error, 1)

	go func(errChan chan<- error) {
		sugar.Infof("Starting client on port: %s\n", cfg.App.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Error(err.Error())
			errChan <- err
		}
	}(errChan)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-quit:
		sugar.Error("Killing signal was received")
	case err := <-errChan:
		sugar.Errorf("HTTP client run error: %s", err)
	}

	sugar.Info("Shutting down client...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.App.AppShutdownTime))
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Infof("Server forced to shutdown: %s", err)
	}
	return nil
}
