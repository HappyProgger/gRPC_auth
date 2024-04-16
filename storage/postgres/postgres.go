package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	cfg "github.com/HappyProgger/gRPC_auth/internal/config"
	"github.com/HappyProgger/gRPC_auth/internal/domain/models"
	"github.com/HappyProgger/gRPC_auth/internal/storage"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// Конструктор Storage
func New(cfg_path string) (*Storage, error) {

	var cfgDB cfg.DbCon
	if err := cleanenv.ReadConfig(cfg_path, &cfgDB); err != nil {
		panic("config path is empty: " + err.Error())
	}
	fmt.Errorf(cfgDB.Password, cfgDB.Username)
	fmt.Errorf(cfgDB.Username, cfgDB.Password,
		cfgDB.DbIP, cfgDB.DbPort, cfgDB.DbName)
	const op = "storage.postgres.New"
	//загружаем конфиг

	//Указываем путь до файла БД

	db, err := sql.Open(`postgres`,
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfgDB.Username, cfgDB.Password, cfgDB.DbIP, cfgDB.DbPort, cfgDB.DbName))
	//fmt.Sprintf("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"
	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	fmt.Errorf("%s: %w", op, "fsdfsadfasdf")

	row := stmt.QueryRowContext(ctx, id)

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

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	// Prepare the query with the RETURNING clause
	var id int64
	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Execute the query, passing parameters
	err = stmt.QueryRowContext(ctx, email, passHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)

		var sqliteErr sqlite3.Error

		// Small hack to identify the ErrConstraintUnique error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Now 'id' contains the ID of the newly inserted row
	return id, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"
	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = $1 ")
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

type Roles struct {
	Admin, User int64
}

var Roles_of_user Roles = Roles{
	Admin: 0,
	User:  1,
}

func (s *Storage) IsAdmin(ctx context.Context, user_id int64) (bool, error) {
	op := "storage.postgres.IsAdmin"
	stmt, err := s.db.Prepare("SELECT role FROM users WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	var role int64
	raw := stmt.QueryRowContext(ctx, user_id)
	err = raw.Scan(&role)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
		//todo сделать нормальную обработку ошибки в случае получения проблемного значения
	}
	//return false, fmt.Errorf("%s: %w", op, fmt.Errorf(reflect.TypeOf(role)))

	//todo что-то не так с проверкой на тип
	return checkRole(role)
}

func checkRole(role int64) (bool, error) {
	switch role {
	case Roles_of_user.Admin:
		return true, nil
	default:
		return false, nil
	}
}
