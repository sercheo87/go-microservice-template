package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-microservice-template/pkg/config"
	"go-microservice-template/pkg/config/logger"
	"go-microservice-template/pkg/controllers"
	"go-microservice-template/pkg/services/configurationService"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func main() {
	config.LoadConfigurationMicroservice("./")
	logger, _ := logger.ApplyLoggerConfiguration(config.Configuration.Log.Level)
	defer closeLoggerHandler()(logger)

	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()
	RequestIDMiddleware()

	logger.Info("Starting server")

	// Commons services instances
	configurationServiceInstance := configurationService.NewConfigurationService()
	serverInstance := controllers.NewServer(configurationServiceInstance)

	stopCh := SetupSignalHandler()

	logger.Info("Creating routers")
	serverInstance.ListenAndServe(stopCh)
}

func closeLoggerHandler() func(logger *zap.Logger) {
	return func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func RequestIDMiddleware() {
	logger.Reset()

	requestID := uuid.New()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestID", requestID)

	logger.LoggerInstance = logger.LoggerInstance.With(zap.String("requestID", requestID.String()))
	logger.WithCtx(ctx, logger.LoggerInstance)
}

var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
