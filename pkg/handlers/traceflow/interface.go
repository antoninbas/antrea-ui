package traceflow

import (
	"context"
)

type RequestsHandler interface {
	CreateRequest(ctx context.Context, request *Request) (string, error)
	GetRequestStatus(ctx context.Context, requestID string) (*RequestStatus, error)
	GetRequestResult(ctx context.Context, requestID string) (map[string]interface{}, error)
	DeleteRequest(ctx context.Context, requestID string) (bool, error)
}
