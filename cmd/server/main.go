package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"antrea.io/antrea-ui/pkg/auth"
	traceflowhandler "antrea.io/antrea-ui/pkg/handlers/traceflow"
	"antrea.io/antrea-ui/pkg/password"
	passwordhasher "antrea.io/antrea-ui/pkg/password/hasher"
	passwordrw "antrea.io/antrea-ui/pkg/password/readwriter"
	"antrea.io/antrea-ui/pkg/server"
	"antrea.io/antrea-ui/pkg/signals"
)

var (
	serverAddr     string
	logger         logr.Logger
	kubeconfig     *string
	privateKeyPath string
)

func run() error {
	var db *sql.DB

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	k8sClient, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		return err
	}

	traceflowHandler := traceflowhandler.NewRequestsHandler(logger, k8sClient)
	passwordStore := password.NewStore(passwordrw.NewInMemory(), passwordhasher.NewArgon2id())
	if err := passwordStore.Init(context.Background()); err != nil {
		return err
	}
	tokenManager := auth.NewTokenManager("key", auth.LoadPrivateKeyOrDie(privateKeyPath))

	s := server.NewServer(logger, db, k8sClient, traceflowHandler, passwordStore, tokenManager)
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
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.StringVar(&privateKeyPath, "private-key", "", "Path to PEM private key file")
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
