package mysql

import (
	"context"

	"github.com/silvan-talos/tlp/example/user"
	"github.com/silvan-talos/tlp/log"
)

// UserRepository is a database mock implementation with customizable behaviour and working happy path.
type UserRepository struct {
	CreateFn func(ctx context.Context, user user.User) (int64, error)
	DeleteFn func(ctx context.Context, id int64) error
	GetFn    func(ctx context.Context, id int64) (*user.User, error)
	UpdateFn func(ctx context.Context, id int64, user user.User) error
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) Create(ctx context.Context, user user.User) (int64, error) {
	log.Debug(ctx, "creating user in database")
	if ur.CreateFn != nil {
		return ur.CreateFn(ctx, user)
	}
	return 1, nil
}

func (ur *UserRepository) Delete(ctx context.Context, id int64) error {
	log.Debug(ctx, "deleting user from database", "id", id)
	if ur.DeleteFn != nil {
		return ur.DeleteFn(ctx, id)
	}
	return nil
}

func (ur *UserRepository) Get(ctx context.Context, id int64) (*user.User, error) {
	log.Debug(ctx, "getting user from database", "id", id)
	if ur.GetFn != nil {
		return ur.GetFn(ctx, id)
	}
	return &user.User{
		ID:   1,
		Name: "Andrei Popescu",
		Age:  53,
	}, nil
}

func (ur *UserRepository) Update(ctx context.Context, id int64, user user.User) error {
	log.Debug(ctx, "updating user in database", "id", id)
	if ur.UpdateFn != nil {
		return ur.UpdateFn(ctx, id, user)
	}
	return nil
}
