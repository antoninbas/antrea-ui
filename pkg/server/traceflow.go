package server

import (
	"github.com/gin-gonic/gin"
)

func (s *server) CreateTraceflowRequest(c *gin.Context) {

}

func (s *server) GetTraceflowRequestStatus(c *gin.Context) {

}

func (s *server) GetTraceflowRequestResult(c *gin.Context) {

}

func (s *server) AddTraceflowRoutes(r *gin.RouterGroup) {
	r = r.Group("/traceflow")
	r.POST("", s.CreateTraceflowRequest)
	r.GET("/:requestId/status", s.GetTraceflowRequestStatus)
	r.GET("/:requestId/result", s.GetTraceflowRequestResult)
}
