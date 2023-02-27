package readwriter

import (
	"context"
	"sync"

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

type K8sSecret struct {
	secretNamespace string
	secretName      string
	k8sClient       dynamic.Interface
}

func (rw *K8sSecret) Read(ctx context.Context) (bool, []byte, []byte, error) {
	return false, nil, nil, nil
}

func (rw *K8sSecret) Write(ctx context.Context, hash []byte, salt []byte) error {
	return nil
}

func NewK8sSecret(secretNamespace string, secretName string, k8sClient dynamic.Interface) *K8sSecret {
	return &K8sSecret{
		secretNamespace: secretNamespace,
		secretName:      secretName,
		k8sClient:       k8sClient,
	}
}
