package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"k8s.io/client-go/dynamic"

	"antrea.io/antrea-ui/pkg/auth"
	traceflowhandler "antrea.io/antrea-ui/pkg/handlers/traceflow"
	"antrea.io/antrea-ui/pkg/password"
)

type server struct {
	logger                   logr.Logger
	db                       *sql.DB
	k8sClient                dynamic.Interface
	traceflowRequestsHandler traceflowhandler.RequestsHandler
	passwordStore            password.Store
	tokenManager             auth.TokenManager
}

func NewServer(
	logger logr.Logger,
	db *sql.DB,
	k8sClient dynamic.Interface,
	traceflowRequestsHandler traceflowhandler.RequestsHandler,
	passwordStore password.Store,
	tokenManager auth.TokenManager,
) *server {
	return &server{
		logger:                   logger,
		db:                       db,
		k8sClient:                k8sClient,
		traceflowRequestsHandler: traceflowRequestsHandler,
		passwordStore:            passwordStore,
		tokenManager:             tokenManager,
	}
}

func tenantUUIDFromURL(c *gin.Context) (string, *serverError) {
	tenantUUID := c.Param("tenantUUID")
	if _, err := uuid.Parse(tenantUUID); err != nil {
		return "", &serverError{
			code:    http.StatusBadRequest,
			message: fmt.Sprintf("Invalid tenantUUID: %s", err.Error()),
		}
	}
	return tenantUUID, nil
}

func (s *server) checkBearerToken(c *gin.Context) {
	if sError := func() *serverError {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			return &serverError{
				code:    http.StatusUnauthorized,
				message: "Missing Authorization header",
			}
		}
		t := strings.Split(auth, " ")
		if len(t) != 2 || t[0] != "Bearer" {
			return &serverError{
				code:    http.StatusUnauthorized,
				message: "Authorization header does not have valid format",
			}
		}
		if err := s.tokenManager.VerifyToken(t[1]); err != nil {
			return &serverError{
				code:    http.StatusUnauthorized,
				message: "Invalid Bearer token",
				err:     err,
			}
		}
		return nil
	}(); sError != nil {
		s.HandleError(c, sError)
		c.Abort()
		return
	}
}

func (s *server) AddRoutes(router *gin.Engine) {
	apiv1 := router.Group("/api/v1")
	s.AddQueryRoutes(apiv1)
	s.AddVariablesRoutes(apiv1)
	s.AddTraceflowRoutes(apiv1)
	s.AddInfoRoutes(apiv1)
	s.AddAccountRoutes(apiv1)
	s.AddAuthRoutes(apiv1)
}
