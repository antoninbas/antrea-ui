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

type Unit string

const (
	UnitBytes   Unit = "Bytes"
	UnitPackets Unit = "Packets"
)

var (
	queriers = map[string]Querier{
		"flowrecords": &FlowRecordsQuerier{},
		"podsentbytes": &PodTrafficQuerier{
			Unit:     UnitBytes,
			Variable: "sourcePod",
			Field:    "source",
		},
		"podsentpackets": &PodTrafficQuerier{
			Unit:     UnitPackets,
			Variable: "sourcePod",
			Field:    "source",
		},
		"podreceivedbytes": &PodTrafficQuerier{
			Unit:     UnitBytes,
			Variable: "destinationPod",
			Field:    "destination",
		},
		"podreceivedpackets": &PodTrafficQuerier{
			Unit:     UnitPackets,
			Variable: "destinationPod",
			Field:    "destination",
		},
		"podnetwork": &PodNetworkQuerier{},
	}
)

type internalQuery struct {
	from       string
	to         string
	intervalMs int32
	variables  []apisv1alpha1.QueryVariable
}

func (q *internalQuery) StringVariable(name string) string {
	for _, v := range q.variables {
		if v.Name == name && v.Type == "string" {
			return v.Value.(string)
		}
	}
	return ""
}

type Querier interface {
	Query(ctx context.Context, db *sql.DB, tenantUUID string, query internalQuery) (*apisv1alpha1.QueryResult, error)
}

type FlowRecordsQuerier struct{}

func (q *FlowRecordsQuerier) Query(
	ctx context.Context,
	db *sql.DB,
	tenantUUID string,
	query internalQuery,
) (*apisv1alpha1.QueryResult, error) {
	// interval := float32(query.intervalMs) / 1000.
	timeValues := make([]interface{}, 0, 1000)
	countValues := make([]interface{}, 0, 1000)
	result := &apisv1alpha1.QueryResult{
		Schema: apisv1alpha1.DataSchema{
			Fields: []apisv1alpha1.SchemaField{
				{
					Name: "time",
					Type: "time",
				},
				{
					Name: "value",
					Type: "int32",
				},
			},
		},
		Values: [][]interface{}{
			timeValues,
			countValues,
		},
	}
	return result, nil
}

type PodTrafficQuerier struct {
	Unit     Unit
	Variable string
	Field    string
}

func (q *PodTrafficQuerier) Query(
	ctx context.Context,
	db *sql.DB,
	tenantUUID string,
	query internalQuery,
) (*apisv1alpha1.QueryResult, error) {
	// interval := float32(query.intervalMs) / 1000.
	timeValues := make([]interface{}, 0, 1000)
	countValues := make([]interface{}, 0, 1000)
	result := &apisv1alpha1.QueryResult{
		Schema: apisv1alpha1.DataSchema{
			Fields: []apisv1alpha1.SchemaField{
				{
					Name: "time",
					Type: "time",
				},
				{
					Name: "value",
					Type: "int64",
				},
			},
		},
		Values: [][]interface{}{
			timeValues,
			countValues,
		},
	}
	return result, nil
}

type PodNetworkQuerier struct{}

func (q *PodNetworkQuerier) Query(
	ctx context.Context,
	db *sql.DB,
	tenantUUID string,
	query internalQuery,
) (*apisv1alpha1.QueryResult, error) {
	sources := make([]interface{}, 0, 1000)
	destinations := make([]interface{}, 0, 1000)
	connectionCounts := make([]interface{}, 0, 1000)
	octetCounts := make([]interface{}, 0, 1000)
	packetCounts := make([]interface{}, 0, 1000)
	reverseOctetCounts := make([]interface{}, 0, 1000)
	reversePacketCounts := make([]interface{}, 0, 1000)
	result := &apisv1alpha1.QueryResult{
		Schema: apisv1alpha1.DataSchema{
			Fields: []apisv1alpha1.SchemaField{
				{
					Name: "source",
					Type: "string",
				},
				{
					Name: "destination",
					Type: "string",
				},
				{
					Name: "connectionCount",
					Type: "int64",
				},
				{
					Name: "octetCount",
					Type: "int64",
				},
				{
					Name: "packetCount",
					Type: "int64",
				},
				{
					Name: "reverseOctetCount",
					Type: "int64",
				},
				{
					Name: "reversePacketCount",
					Type: "int64",
				},
			},
		},
		Values: [][]interface{}{
			sources,
			destinations,
			connectionCounts,
			octetCounts,
			packetCounts,
			reverseOctetCounts,
			reversePacketCounts,
		},
	}
	return result, nil
}

func bindQueryJSON(c *gin.Context, query *apisv1alpha1.Query) error {
	if err := c.BindJSON(query); err != nil {
		return err
	}
	return nil
}

func (s *server) RunQuery(c *gin.Context) {
	query := apisv1alpha1.Query{}
	if sError := func() *serverError {
		tenantUUID := tenantUUIDFromContext(c)
		if err := bindQueryJSON(c, &query); err != nil {
			return &serverError{
				code:    http.StatusBadRequest,
				message: err.Error(),
			}
		}
		queryName := query.QueryName
		querier, ok := queriers[queryName]
		if !ok {
			return &serverError{
				code:    http.StatusBadRequest,
				message: fmt.Sprintf("No matching querier for query name '%s'", queryName),
			}
		}
		now := time.Now()
		// by default, go 1 year back
		from, err := timestamps.ParseTimestamp(query.From, now, now.AddDate(-1, 0, 0))
		if err != nil {
			return &serverError{
				code:    http.StatusBadRequest,
				message: err.Error(),
			}
		}
		to, err := timestamps.ParseTimestamp(query.To, now, now)
		if err != nil {
			return &serverError{
				code:    http.StatusBadRequest,
				message: err.Error(),
			}
		}
		timeoutMs := 60000 // 1 minute
		if query.TimeoutMs > 0 {
			timeoutMs = int(query.TimeoutMs)
		}
		ctx, cancel := context.WithTimeout(c, time.Duration(timeoutMs)*time.Millisecond)
		defer cancel()
		result, err := querier.Query(
			ctx,
			s.db,
			tenantUUID,
			internalQuery{
				from:       from,
				to:         to,
				intervalMs: query.IntervalMs,
				variables:  query.Variables,
			},
		)
		if err != nil {
			return &serverError{
				code: http.StatusInternalServerError,
				err:  err,
			}
		}
		result.RefID = query.RefID
		c.JSON(http.StatusOK, result)
		return nil
	}(); sError != nil {
		s.HandleError(c, sError)
		s.LogError(sError, "Failed to run Query", "refID", query.RefID)
		return
	}
}

func (s *server) AddQueryRoutes(r *gin.RouterGroup) {
	r = r.Group("/query")
	r.POST("", s.RunQuery)
}
