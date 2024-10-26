package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	jwt_lib "github.com/babs-corp/babs-maps-auth/internal/lib/jwt"
	"github.com/babs-corp/babs-maps-auth/internal/lib/logger/handlers/sl"
	"github.com/babs-corp/babs-maps-auth/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
	secret       string
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid uuid.UUID, err error)
}

// We can get user not only from Database, but e.g. from kafka, cache, etc...
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	Users(ctx context.Context, limit uint) ([]models.User, error)
	UserById(ctx context.Context, id uuid.UUID) (models.User, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new instance of Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
	secret string,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
		secret:       secret,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid password", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// app, err := a.appProvider.App(ctx, appId)
	// if err != nil {
	// 	if errors.Is(err, storage.ErrAppNotFound) {
	// 		a.log.Warn("app not found", sl.Err(err))
	// 		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	// 	}
	// 	a.log.Error("failed to get app", sl.Err(err))
	// }
	// a.log.Info("app found", slog.String("app", app.Name))

	token, err := jwt_lib.NewToken(user, a.secret, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID
// If user already exists, returns error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (uuid.UUID, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return uuid.UUID{}, fmt.Errorf("%w", ErrUserExists)
		}
		log.Error("failed to save user", sl.Err(err))

		return uuid.UUID{}, fmt.Errorf("cannot register user")
	}

	return id, nil
}

// IsAdmin checks if user is admin
func (a *Auth) IsAdmin(
	ctx context.Context,
	userId uuid.UUID,
) (bool, error) {
	const op = "auth.IsAdmin"

	a.log.With(
		slog.String("op", op),
		slog.String("user_id", userId.String()),
	)

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		a.log.Error("failed to check if user is admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (a *Auth) UserById(
	ctx context.Context,
	userId uuid.UUID,
) (models.User, error) {
	const op = "auth.UserById"

	a.log.With(
		slog.String("op", op),
		slog.String("user_id", userId.String()),
	)

	user, err := a.userProvider.UserById(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		a.log.Error("failed to check if user is admin", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (a *Auth) Users(
	ctx context.Context,
	limit uint,
) ([]models.User, error) {
	const op = "auth.Users"

	a.log.With(
		slog.String("op", op),
	)

	users, err := a.userProvider.Users(ctx, limit)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return []models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		a.log.Error("failed to check if user is admin", sl.Err(err))
		return []models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (a *Auth) ValidateToken(
	ctx context.Context,
	token string,
) (uuid.UUID, error) {
	const op = "auth.validateToken"

	a.log.With(
		slog.String("op", op),
	)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[%s] cannot validate token: %w", op, err)
	}

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("[%s] cannot parse token claims", op)
	}

	str_uid, ok := claims["uid"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("[%s] cannot parse string uuid", op)
	}
	uid, err := uuid.Parse(str_uid)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[%s] cannot parse uuid claims", op)
	}
	
	return uid, nil
}
