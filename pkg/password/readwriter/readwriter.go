package readwriter

import (
	"context"
	"fmt"
	"sync"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Interface interface {
	Read(ctx context.Context) (bool, []byte, []byte, error)
	Write(ctx context.Context, hash []byte, salt []byte) error
}

type InMemory struct {
	sync.Mutex
	set  bool
	hash []byte
	salt []byte
}

func (rw *InMemory) Read(ctx context.Context) (bool, []byte, []byte, error) {
	rw.Lock()
	defer rw.Unlock()
	if !rw.set {
		return false, nil, nil, nil
	}
	return true, rw.hash, rw.salt, nil
}

func (rw *InMemory) Write(ctx context.Context, hash []byte, salt []byte) error {
	rw.Lock()
	defer rw.Unlock()
	rw.hash = hash
	rw.salt = salt
	rw.set = true
	return nil
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

var (
	k8sSecretGVR = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}
)

type K8sSecret struct {
	secretNamespace string
	secretName      string
	k8sClient       dynamic.Interface
}

func (rw *K8sSecret) readSecret(ctx context.Context) (*unstructured.Unstructured, bool, error) {
	secret, err := rw.k8sClient.Resource(k8sSecretGVR).Namespace(rw.secretNamespace).Get(ctx, rw.secretName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return secret, true, nil
}

func (rw *K8sSecret) Read(ctx context.Context) (bool, []byte, []byte, error) {
	secret, ok, err := rw.readSecret(ctx)
	if !ok {
		return false, nil, nil, nil
	}
	if err != nil {
		return false, nil, nil, fmt.Errorf("error when retrieving K8s secret '%s/%s': %w", rw.secretNamespace, rw.secretName, err)
	}

	readData := func() ([]byte, []byte, error) {
		data, ok, err := unstructured.NestedMap(secret.Object, "data")
		if err != nil {
			// should not be possible
			return nil, nil, err
		}
		if !ok {
			// should not be possible
			return nil, nil, fmt.Errorf("no data in secret")
		}
		hash, ok := data["hash"]
		if !ok {
			return nil, nil, fmt.Errorf("hash is missing from data")
		}
		salt, ok := data["salt"]
		if !ok {
			return nil, nil, fmt.Errorf("salt is missing from data")
		}
		return hash.([]byte), salt.([]byte), nil
	}

	hash, salt, err := readData()
	if err != nil {
		return false, nil, nil, fmt.Errorf("error when reading data from K8s secret '%s/%s': %w", rw.secretNamespace, rw.secretName, err)
	}

	return true, hash, salt, nil
}

func (rw *K8sSecret) Write(ctx context.Context, hash []byte, salt []byte) error {
	secret, ok, err := rw.readSecret(ctx)
	if err != nil {
		return err
	}
	if !ok {
		// create
		secret := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": k8sSecretGVR.Group + "/" + k8sSecretGVR.Version,
				"kind":       "Secret",
				"metadata": map[string]interface{}{
					"namespace": rw.secretNamespace,
					"name":      rw.secretName,
				},
				"data": map[string][]byte{
					"hash": hash,
					"salt": salt,
				},
			},
		}
		if _, err := rw.k8sClient.Resource(k8sSecretGVR).Namespace(rw.secretNamespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("error when creating K8s secret '%s/%s': %w", rw.secretNamespace, rw.secretName, err)
		}
		return nil
	}
	// update
	if _, err := rw.k8sClient.Resource(k8sSecretGVR).Namespace(rw.secretNamespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
		// we do not handle update conflicts, as we should be the only writer
		return fmt.Errorf("error when updating K8s secret '%s/%s': %w", rw.secretNamespace, rw.secretName, err)
	}
	return nil
}

func NewK8sSecret(secretNamespace string, secretName string, k8sClient dynamic.Interface) *K8sSecret {
	return &K8sSecret{
		secretNamespace: secretNamespace,
		secretName:      secretName,
		k8sClient:       k8sClient,
	}
}
