package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/HappyProgger/gRPC_auth/internal/domain/models"
	"github.com/HappyProgger/gRPC_auth/internal/lib/jwt"
	"github.com/HappyProgger/gRPC_auth/internal/lib/logger/sl"
	"github.com/HappyProgger/gRPC_auth/internal/storage"
	"github.com/HappyProgger/gRPC_auth/storage/postgres"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider *postgres.Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL, // Время жизни возвращаемых токенов
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	op := "Auth.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("register new user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userId, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("user already exist", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return userId, nil

}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *Auth) Login(ctx context.Context,
	email string,
	password string, // пароль в чистом виде, аккуратней с логами!
	appID int,       // ID приложения, в котором логинится пользователь
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("username", email)) // password либо не логируем, либо логируем в замаскированном виде

	log.Info("attempting to login user")

	// Достаем пользователя из БД
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем корректность полученного пароля
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	//todo доделать
	// Получаем информацию о приложении
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	// Создаем токен авторизации
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
