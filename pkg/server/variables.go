package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apisv1alpha1 "antrea.io/antrea-ui/apis/v1alpha1"
	"antrea.io/antrea-ui/pkg/utils/timestamps"
)

var (
	listers = map[string]VariableValueLister{
		"sourcePod":      &PodLister{Field: "source"},
		"destinationPod": &PodLister{Field: "destination"},
	}
)

type VariableValueLister interface {
	List(ctx context.Context, db *sql.DB, name, from, to string) ([]interface{}, string, error)
}

type PodLister struct {
	Field string
}

func (l *PodLister) List(ctx context.Context, db *sql.DB, name string, from string, to string) ([]interface{}, string, error) {
	values := make([]interface{}, 0, 100)
	return values, "string", nil
}

func (s *server) GetVariables(c *gin.Context) {
	variables := []apisv1alpha1.Variable{
		{
			Name: "sourcePod",
			Type: "string",
		},
		{
			Name: "destinationPod",
			Type: "string",
		},
	}
	c.JSON(http.StatusOK, variables)
}

func (s *server) GetVariable(c *gin.Context) {
	name := c.Param("name")
	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")
	if sError := func() *serverError {
		lister, ok := listers[name]
		if !ok {
			return &serverError{
				code:    http.StatusBadRequest,
				message: fmt.Sprintf("Unknown variable name '%s'", name),
			}
		}
		now := time.Now()
		// by default, go 1 year back
		from, err := timestamps.ParseTimestamp(from, now, now.AddDate(-1, 0, 0))
		if err != nil {
			return &serverError{
				code:    http.StatusBadRequest,
				message: err.Error(),
			}
		}
		to, err := timestamps.ParseTimestamp(to, now, now)
		if err != nil {
			return &serverError{
				code:    http.StatusBadRequest,
				message: err.Error(),
			}
		}
		values, variableType, err := lister.List(c, s.db, name, from, to)
		if err != nil {
			return &serverError{
				code: http.StatusInternalServerError,
				err:  err,
			}
		}
		variable := apisv1alpha1.Variable{
			Name:   name,
			Type:   variableType,
			Values: values,
		}
		c.JSON(http.StatusOK, variable)
		return nil
	}(); sError != nil {
		s.HandleError(c, sError)
		s.LogError(sError, "Failed to get variable info", "name", name)
		return
	}
}

func (s *server) AddVariablesRoutes(r *gin.RouterGroup) {
	r = r.Group("/variables")
	r.GET("", s.GetVariables)
	r.GET("/:name", s.GetVariable)
}
