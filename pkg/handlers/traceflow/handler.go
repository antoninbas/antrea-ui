package traceflow

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
)

const (
	traceflowExpiryTimeout = 60 * time.Minute
)

var (
	traceflowGVR = schema.GroupVersionResource{
		Group:    "crd.antrea.io",
		Version:  "v1alpha1",
		Resource: "traceflows",
	}
)

type requestsHandler struct {
	logger    logr.Logger
	k8sClient dynamic.Interface
}

func NewRequestsHandler(logger logr.Logger, k8sClient dynamic.Interface) *requestsHandler {
	return &requestsHandler{
		logger:    logger,
		k8sClient: k8sClient,
	}
}

func (h *requestsHandler) Run(stopCh <-chan struct{}) {
	go h.runGC(stopCh)
	<-stopCh
}

func (h *requestsHandler) CreateRequest(ctx context.Context, request *Request) (string, error) {
	requestID := uuid.NewString()
	if err := h.createTraceflow(ctx, requestID, request.Object); err != nil {
		return "", err
	}
	return requestID, nil
}

func (h *requestsHandler) GetRequestStatus(ctx context.Context, requestID string) (*RequestStatus, error) {
	_, done, err := h.getTraceflow(ctx, requestID)
	return &RequestStatus{
		Done: done,
		Err:  err,
	}, nil
}

func (h *requestsHandler) GetRequestResult(ctx context.Context, requestID string) (map[string]interface{}, error) {
	object, done, err := h.getTraceflow(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if !done {
		return nil, fmt.Errorf("Traceflow not complete yet")
	}
	return object, nil
}

func (h *requestsHandler) getTraceflow(ctx context.Context, tfName string) (map[string]interface{}, bool, error) {
	traceflow, err := h.k8sClient.Resource(traceflowGVR).Get(ctx, tfName, metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}
	phase, ok, err := unstructured.NestedString(traceflow.Object, "status", "phase")
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return traceflow.Object, false, nil
	}
	return traceflow.Object, (phase == "Succeeded" || phase == "Failed"), nil
}

func (h *requestsHandler) createTraceflow(ctx context.Context, tfName string, object map[string]interface{}) error {
	traceflow := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": traceflowGVR.Group + "/" + traceflowGVR.Version,
			"kind":       "Traceflow",
			"metadata": map[string]interface{}{
				"name": tfName,
			},
			"spec": object["spec"],
		},
	}
	if _, err := h.k8sClient.Resource(traceflowGVR).Create(ctx, traceflow, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func (h *requestsHandler) doGC(ctx context.Context) {
	list, err := h.k8sClient.Resource(traceflowGVR).List(ctx, metav1.ListOptions{})
	if err != nil {
		h.logger.Error(err, "Error when listing traceflows")
		return
	}
	expiredTraceflows := []string{}
	now := time.Now()
	for _, tf := range list.Items {
		creationTimestamp := tf.GetCreationTimestamp()
		if now.Sub(creationTimestamp.Time) > traceflowExpiryTimeout {
			expiredTraceflows = append(expiredTraceflows, tf.GetName())
		}
	}
	for _, tfName := range expiredTraceflows {
		if err := h.k8sClient.Resource(traceflowGVR).Delete(ctx, tfName, metav1.DeleteOptions{}); err != nil {
			h.logger.Error(err, "Error when deleting expired traceflow", "name", tfName)
		}
	}
}

func (h *requestsHandler) runGC(stopCh <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wait.UntilWithContext(ctx, h.doGC, 1*time.Minute)
	<-stopCh
}
