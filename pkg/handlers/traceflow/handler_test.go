package traceflow

// import (
// 	"testing"
// 	"time"

// 	"github.com/go-logr/zapr"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/zap"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
// 	"k8s.io/apimachinery/pkg/util/wait"
// 	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// 	"sigs.k8s.io/controller-runtime/pkg/client/fake"
// )

// var (
// 	scheme = runtime.NewScheme()
// )

// func init() {
// 	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
// }

// type testData struct {
// 	*testing.T
// 	stopCh    chan struct{}
// 	handler   *requestsHandler
// 	k8sClient client.Client
// }

// func setUp(t *testing.T, objs ...client.Object) *testData {
// 	zc := zap.NewDevelopmentConfig()
// 	zapLog, err := zc.Build()
// 	require.NoError(t, err)
// 	logger := zapr.NewLogger(zapLog)
// 	k8sClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
// 	handler := NewRequestsHandler(logger, k8sClient)
// 	stopCh := make(chan struct{})
// 	go handler.Run(stopCh)
// 	return &testData{
// 		T:         t,
// 		stopCh:    stopCh,
// 		handler:   handler,
// 		k8sClient: k8sClient,
// 	}
// }

// func tearDown(t *testData) {
// 	close(t.stopCh)
// }

// func TestRequestsHandler(t *testing.T) {
// 	traceflowNamespace := "default"
// 	traceflowName := "my-traceflow"
// 	role := "view"
// 	traceflow := &corev1.Traceflow{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Namespace: traceflowNamespace,
// 			Name:      traceflowName,
// 		},
// 		Data: map[string][]byte{
// 			"kubeconfig": []byte("TOP TRACEFLOW STUFF"),
// 		},
// 	}
// 	testData := setUp(t, traceflow)
// 	handler := testData.handler
// 	requestID, err := handler.EnqueueRequest(&Request{
// 		Namespace: traceflowNamespace,
// 		Name:      traceflowName,
// 		Role:      role,
// 	})
// 	require.NoError(t, err)
// 	var requestStatus *RequestStatus
// 	err = wait.Poll(100*time.Millisecond, 1*time.Second, func() (bool, error) {
// 		var err error
// 		requestStatus, err = handler.GetRequestStatus(requestID)
// 		if err != nil {
// 			return false, err
// 		}
// 		return requestStatus.Done, nil
// 	})
// 	require.NoError(t, err)
// 	require.NoError(t, requestStatus.Err)
// 	traceflowID := requestStatus.TraceflowID
// 	require.NotEmpty(t, traceflowID)
// 	connectionTraceflow, err := handler.GetTraceflow(traceflowID)
// 	require.NoError(t, err)
// 	assert.Equal(t, role, connectionTraceflow.Role)
// 	assert.Equal(t, traceflow.Data["kubeconfig"], connectionTraceflow.Data)
// }
