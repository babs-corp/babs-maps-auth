package storage

import (
	"context"
	"errors"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type Storage interface {
	User(ctx context.Context, email string) (*models.User, error)
	App(ctx context.Context, appId int) (*models.App, error)
}
