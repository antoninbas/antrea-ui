package password

import (
	"context"
	"crypto/rand"
	"fmt"
	"golang.org/x/exp/slices"
	"sync"

	"antrea.io/antrea-ui/pkg/password/hasher"
	"antrea.io/antrea-ui/pkg/password/readwriter"
)

const defaultPassword = "admin"

var (
	NotInitializedErr  = fmt.Errorf("not initialized")
	InvalidPasswordErr = fmt.Errorf("invalid password")
)

type Store interface {
	Update(ctx context.Context, password []byte) error
	Compare(ctx context.Context, password []byte) error
}

type store struct {
	sync.Mutex
	cachedSalt []byte
	cachedHash []byte
	rw         readwriter.Interface
	hasher     hasher.Interface
}

func NewStore(rw readwriter.Interface, hasher hasher.Interface) *store {
	return &store{
		rw:     rw,
		hasher: hasher,
	}
}

func (s *store) Init(ctx context.Context) error {
	ok, hash, salt, err := s.rw.Read(ctx)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	if ok {
		s.cachedSalt = salt
		s.cachedHash = hash
		return nil
	}
	salt = make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("error when generation random salt: %w", err)
	}
	hash, err = s.hasher.Hash([]byte(defaultPassword), salt)
	if err != nil {
		return err
	}
	if err := s.rw.Write(ctx, hash, salt); err != nil {
		return err
	}
	s.cachedSalt = salt
	s.cachedHash = hash
	return nil
}

func (s *store) Update(ctx context.Context, password []byte) error {
	s.Lock()
	defer s.Unlock()
	if s.cachedSalt == nil {
		return NotInitializedErr
	}
	hash, err := s.hasher.Hash(password, s.cachedSalt)
	if err != nil {
		return err
	}
	s.cachedHash = hash
	return nil
}

func (s *store) Compare(ctx context.Context, password []byte) error {
	// TODO: better lock granularity
	s.Lock()
	defer s.Unlock()
	if s.cachedSalt == nil {
		return NotInitializedErr
	}
	hash, err := s.hasher.Hash(password, s.cachedSalt)
	if err != nil {
		return err
	}
	if !slices.Equal(hash, s.cachedHash) {
		return InvalidPasswordErr
	}
	return nil
}
