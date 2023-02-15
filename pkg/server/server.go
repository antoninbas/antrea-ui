package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"

	traceflowhandler "antrea.io/antrea-ui/pkg/handlers/traceflow"
)

type server struct {
	logger                   logr.Logger
	db                       *sql.DB
	traceflowRequestsHandler traceflowhandler.RequestsHandler
}

func NewServer(
	logger logr.Logger,
	db *sql.DB,
	traceflowRequestsHandler traceflowhandler.RequestsHandler,
) *server {
	return &server{
		logger:                   logger,
		db:                       db,
		traceflowRequestsHandler: traceflowRequestsHandler,
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

func (s *server) AddRoutes(router *gin.Engine) {
	apiv1 := router.Group("/api/v1")
	// apiv1.Use(func(c *gin.Context) {
	// 	tenantUUID, sError := tenantUUIDFromURL(c)
	// 	if sError != nil {
	// 		s.HandleError(c, sError)
	// 		c.Abort()
	// 		return
	// 	}
	// 	c.Set("tenantUUID", tenantUUID)
	// })
	s.AddQueryRoutes(apiv1)
	s.AddVariablesRoutes(apiv1)
	s.AddTraceflowRoutes(apiv1)
}
