package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	"github.com/babs-corp/babs-maps-auth/internal/services/auth"
	"github.com/babs-corp/babs-maps-auth/internal/storage"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var _ auth.UserProvider = (*Storage)(nil)
var _ auth.UserSaver = (*Storage)(nil)
var _ auth.AppProvider = (*Storage)(nil)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite"
	fmt.Println("storage path: ", storagePath)
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passwordHash []byte) (userId uuid.UUID, err error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?) RETURNING id")
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res := stmt.QueryRowContext(ctx, email, passwordHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.Is(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	var id uuid.UUID
	err = res.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "storage.sqlite.UserById"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE id = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.sqlite.App"
	stmt, err := s.db.Prepare("SELECT * FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
