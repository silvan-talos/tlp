package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/silvan-talos/tlp/example"
	"github.com/silvan-talos/tlp/log"
)

type User struct {
	ID   int64
	Name string
	Age  int
}

type Repository interface {
	Create(ctx context.Context, user User) (int64, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, id int64, user User) error
}

type Service struct {
	users Repository
}

func NewService(users Repository) *Service {
	return &Service{
		users: users,
	}
}

func (s *Service) CreateUser(ctx context.Context, user User) (int64, error) {
	id, err := s.users.Create(ctx, user)
	if err != nil {
		log.Error(ctx, "create user", "err", err, "user", user)
		return 0, errors.New("failed to create user")
	}
	log.Info(ctx, "user created successfully", "id", id)
	return id, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.users.Get(ctx, id)
	if err != nil {
		log.Error(ctx, "get user to delete", "err", err, "id", id)
		if errors.Is(err, example.ErrNotFound) {
			// fail silently
			return nil
		}
		return example.ErrInternal
	}
	err = s.users.Delete(ctx, id)
	if err != nil {
		log.Error(ctx, "delete user", "err", err, "id", id)
		return errors.New("failed to delete user")
	}
	log.Info(ctx, "user deleted successfully", "id", id)
	return nil
}

func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
	user, err := s.users.Get(ctx, id)
	if err != nil {
		log.Error(ctx, "get user", "err", err, "id", id)
		if errors.Is(err, example.ErrNotFound) {
			return nil, fmt.Errorf("user %w", example.ErrNotFound)
		}
		return nil, example.ErrInternal
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int64, user User) error {
	current, err := s.users.Get(ctx, id)
	if err != nil {
		log.Error(ctx, "get user to update", "err", err, "id", id)
		return fmt.Errorf("user with id %d can't be updated", id)
	}
	if user.Name != "" && user.Name != current.Name {
		current.Name = user.Name
	}
	if user.Age > 0 && user.Age != current.Age {
		current.Age = user.Age
	}
	err = s.users.Update(ctx, id, *current)
	if err != nil {
		log.Error(ctx, "update user", "err", err, "id", id)
		return errors.New("failed to update user")
	}
	log.Info(ctx, "user updated successfully", "id", id)
	return nil
}
