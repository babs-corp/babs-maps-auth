package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	"github.com/babs-corp/babs-maps-auth/internal/services/auth"
	"github.com/babs-corp/babs-maps-auth/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var _ auth.UserProvider = (*Storage)(nil)
var _ auth.UserSaver = (*Storage)(nil)
var _ auth.AppProvider = (*Storage)(nil)

type Storage struct {
	db *sqlx.DB
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.pgx"

	db, err := sqlx.Connect("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passwordHash []byte) (userId uuid.UUID, err error) {
	const op = "storage.pgx.SaveUser"

	query := fmt.Sprintf("INSERT INTO users (email, pass_hash) VALUES ('%s', '%s') RETURNING id;", email, string(passwordHash))
	stmt, err := s.db.PreparexContext(ctx, query)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res := stmt.QueryRowxContext(ctx)
	var id uuid.UUID
	err = res.Scan(&id)
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == UniqueViolation {
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.pgx.User"

	query := fmt.Sprintf("SELECT id, email, pass_hash FROM users WHERE email = '%s'", email)
	stmt, err := s.db.PreparexContext(ctx, query)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx)

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

func (s *Storage) UserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.pgx.UserById"

	query := fmt.Sprintf("SELECT id, email, pass_hash, created_at FROM users WHERE id = '%s'", uid)
	stmt, err := s.db.PreparexContext(ctx, query)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx)

	var user models.User
	err = row.StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) Users(ctx context.Context, limit uint) ([]models.User, error) {
	const op = "storage.pgx.Users"

	query := fmt.Sprintf("SELECT * FROM users LIMIT '%d'", limit)
	stmt, err := s.db.PreparexContext(ctx, query)
	if err != nil {
		return []models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return []models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.StructScan(&user)
		if err != nil {
			return []models.User{}, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "storage.pfx.IsAdmin"

	stmt, err := s.db.PreparexContext(ctx, "SELECT is_admin FROM users WHERE id = ?")
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
	const op = "storage.pgx.App"
	stmt, err := s.db.PreparexContext(ctx, "SELECT * FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, appID)

	var app models.App
	err = row.StructScan(&app)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
