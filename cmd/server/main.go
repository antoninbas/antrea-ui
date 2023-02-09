package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"

	"antrea.io/antrea-ui/pkg/server"
	"antrea.io/antrea-ui/pkg/signals"
)

var (
	serverAddr string
	logger     logr.Logger
)

func run() error {
	var db *sql.DB
	s := server.NewServer(logger, db)
	router := gin.Default()
	// TODO(antonin): CHANGEME
	// This is for testing, it should be configurable
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))
	s.AddRoutes(router)

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	stopCh := signals.RegisterSignalHandlers()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logger.Info("Starting server", "address", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err, "Server error")
			os.Exit(1)
		}
	}()

	<-stopCh

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	return nil
}

func main() {
	flag.StringVar(&serverAddr, "addr", ":8080", "Listening address for server")
	flag.Parse()

	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zc.DisableStacktrace = true
	zapLog, err := zc.Build()
	if err != nil {
		panic("Cannot initialize Zap logger")
	}
	logger = zapr.NewLogger(zapLog)
	if err := run(); err != nil {
		logger.Error(err, "error in run() function")
		os.Exit(1)
	}
}
