package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apisv1alpha1 "antrea.io/antrea-ui/apis/v1alpha1"
)

func (s *server) GetToken(c *gin.Context) {
	if sError := func() *serverError {
		user, password, ok := c.Request.BasicAuth()
		if !ok {
			return &serverError{
				code:    http.StatusUnauthorized,
				message: "Basic Auth required",
			}
		}
		if user != "admin" {
			return &serverError{
				code:    http.StatusUnauthorized,
				message: "Must authenticate as admin",
			}
		}
		if err := s.passwordStore.Compare(c, []byte(password)); err != nil {
			return &serverError{
				code:    http.StatusUnauthorized,
				message: "Invalid admin password",
			}
		}
		token, err := s.tokenManager.GetToken()
		if err != nil {
			return &serverError{
				code: http.StatusInternalServerError,
				err:  fmt.Errorf("error when getting JWT token: %w", err),
			}
		}
		resp := apisv1alpha1.GetToken{
			AccessToken: token.Raw,
			TokenType:   "Bearer",
			ExpiresIn:   int64(token.ExpiresIn / time.Second),
		}
		c.JSON(http.StatusOK, resp)
		return nil
	}(); sError != nil {
		s.HandleError(c, sError)
		s.LogError(sError, "Failed to get token")
		return
	}
}

func (s *server) AddAuthRoutes(r *gin.RouterGroup) {
	r = r.Group("/auth")
	r.GET("/login", s.GetToken)
}
