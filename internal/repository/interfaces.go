package repository

import (
	"context"

	"github.com/ssr0016/template/internal/model"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	CreateWithPassword(ctx context.Context, email, name, hash string) (*model.User, error)
	List(ctx context.Context, emailFilter string, limit int) ([]model.User, error)
}

// Compile-time check: UserRepo implements UserRepository.
var _ UserRepository = (*UserRepo)(nil)
