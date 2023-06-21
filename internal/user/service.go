package user

import (
	"context"
	"github.com/altiby/son/internal/domain"
	"github.com/google/uuid"
)

type Storage interface {
	RegisterUser(ctx context.Context, user domain.User, password string) error
	AuthorizeUser(ctx context.Context, id string, password string) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	SearchUsers(ctx context.Context, firstName string, lastName string) ([]domain.User, error)
}

type Hasher interface {
	Hash(in string) (string, error)
}

type Service struct {
	storage Storage
	hasher  Hasher
}

func (s Service) SearchUsers(ctx context.Context, firstName string, lastName string) ([]domain.User, error) {
	return s.storage.SearchUsers(ctx, firstName, lastName)
}

func (s Service) RegisterUser(ctx context.Context, user domain.User, password string) (domain.User, error) {
	user.ID = uuid.NewString()
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return domain.User{}, err
	}
	return user, s.storage.RegisterUser(ctx, user, passwordHash)
}

func (s Service) AuthorizeUser(ctx context.Context, id string, password string) (domain.User, error) {
	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return domain.User{}, err
	}
	return s.storage.AuthorizeUser(ctx, id, passwordHash)
}

func (s Service) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return s.storage.GetUserByID(ctx, id)
}

func NewService(storage Storage, hasher Hasher) Service {
	return Service{storage: storage, hasher: hasher}
}
