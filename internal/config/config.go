package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	CfgPath        string     `yaml:"config_path" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	DbCon          DbCon         `yaml:"db"`
	Clients        ClientConfig  `yaml:"client_config"`
	AppSecret      string        `yaml:"app_secret" env:"APP_SECRET"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DbCon struct {
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	DbName   string `yaml:"db_name" env-default:"postgres"`
	DbIP     string `yaml:"db_ip" env-default:"localhost"`
	DbPort   string `yaml:"db_port" env-default:"5432"`
}

type Client struct {
	Addres       string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}
type ClientConfig struct {
	SSO Client `yaml:"sso"`
}

var Cfg Config

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	var Cfg Config
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &Cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &Cfg
}

func MustLoadWithPath(pathToConfig string) *Config {
	configPath := pathToConfig
	if pathToConfig == "" {
		configPath = fetchConfigPath()
	}
	var Cfg Config

	if err := cleanenv.ReadConfig(configPath, &Cfg); err != nil {
		panic("config path is invalid: " + err.Error())
	}

	return &Cfg
}

func configPath() string {
	conf := "CONFIG_PATH"
	if os.Getenv(conf); os.Getenv(conf) != "" {
		return os.Getenv(conf)
	}

	return "./internal/config/config_local.yaml"

}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

//Создаем логгер с указанием текущего окружения.
//
//envLocal — локальный запуск. Используем удобный для консоли TextHandler и уровень логиро
