package controllers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go-microservice-template/pkg/config"
	"go-microservice-template/pkg/config/logger"
	"go-microservice-template/pkg/services/configurationService"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Router               *mux.Router
	Handler              http.Handler
	ConfigurationService configurationService.Executor
}

func NewServer(configurationService configurationService.Executor) *Server {
	return &Server{
		Router:               mux.NewRouter(),
		ConfigurationService: configurationService,
	}
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	s.registerHandlers()
	s.Handler = s.Router
	logger.FromCtx().Info("Starting server...")

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Configuration.MicroserviceServer, config.Configuration.MicroservicePort), RequestIDMiddleware(s.Handler)); err != nil {
			logger.FromCtx().Fatal("Error starting server", zap.Error(err))
		}
	}()
	logger.FromCtx().Info("Server started")

	// wait for SIGTERM or SIGINT
	<-stopCh
	logger.FromCtx().Info("Stopping server...")
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Reset()

		requestID := uuid.New()

		ctx := r.Context()
		const keyTraceIdContext = "trace-id"
		ctx = context.WithValue(ctx, keyTraceIdContext, requestID)

		r = r.WithContext(ctx)

		logger.LoggerInstance = logger.LoggerInstance.With(zap.String(keyTraceIdContext, requestID.String()))
		r = r.WithContext(logger.WithCtx(ctx, logger.LoggerInstance))

		next.ServeHTTP(w, r)
	})
}

func (s *Server) registerHandlers() {
	s.Router.HandleFunc("/api/hello",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello World"))
		},
	).Methods("GET")
}
