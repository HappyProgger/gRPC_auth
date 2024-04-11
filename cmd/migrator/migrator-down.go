package main

//
//import (
//	"errors"
//	"flag"
//	"fmt"
//	// Библиотека для миграций
//	"github.com/golang-migrate/migrate/v4"
//	// Драйвер для выполнения миграций SQLite 3
//	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
//	// Драйвер для получения миграций из файлов
//	_ "github.com/golang-migrate/migrate/v4/source/file"
//)
//
//func main() {
//	var storagePath, migrationsPath, migrationsTable string
//
//	// Получаем необходимые значения из флагов запуска
//
//	// Путь до файла БД
//	// Его достаточно, т.к. мы используем SQLite, другие креды не нужны
//	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
//	// Путь до папки с миграциями
//	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
//	// Таблица, в которой будет храниться информация о миграциях. Она нужна
//	// для того, чтобы понимать, какие миграции уже применены, а какие нет
//	// Дефолтное значение — 'migrations'
//	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
//	flag.Parse()
//	// Выполняем парсинг флагов
//
//	// Валидация параметров
//	if storagePath == "" {
//		// Простейший способ обработки ошибки
//		// При необходимости можете выбрать более подходящий вариант
//		// Меня устраивает паника, поскольку это вспомогательная утилита
//		panic("storage-path is required")
//	}
//	if migrationsPath == "" {
//		panic("migrations-path is required")
//	}
//
//	// Создаем объект мигратора, передав креды нашей БД
//	m, err := migrate.New("file://"+migrationsPath,
//		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	// Выполняем миграции до последней версии
//	if err := m.Up(); err != nil {
//		if errors.Is(err, migrate.ErrNoChange) {
//			fmt.Println("no migrations to apply")
//			return
//		}
//		panic(err)
//
//	}
//	fmt.Println("migration was successful")
//
//}

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	var username, password, dbName, dbIP, dbPort, storagePath, migrationsPath, migrationsTable string

	// Получаем необходимые значения из флагов запуска

	// Путь до файла БД
	// Его достаточно, т.к. мы используем SQLite, другие креды не нужны

	flag.StringVar(&username, "username", "postgres", "username")
	flag.StringVar(&password, "password", "postgres", "password from DB")
	flag.StringVar(&dbName, "dbName", "postgres", "name of db")
	flag.StringVar(&dbIP, "dbIP", "localhost", "db ip")
	flag.StringVar(&dbPort, "dbPort", "5432", "db port")

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	// Путь до папки с миграциями
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	// Таблица, в которой будет храниться информация о миграциях. Она нужна
	// для того, чтобы понимать, какие миграции уже применены, а какие нет
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()
	// Выполняем парсинг флагов

	// Валидация параметров
	if storagePath == "" {
		// Простейший способ обработки ошибки
		// При необходимости можете выбрать более подходящий вариант
		// Меня устраивает паника, поскольку это вспомогательная утилита
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	db, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, dbIP, dbPort, dbName),
	)

	// ex config "postgres://username:password@localhost:5432/database_name?sslmode=disable"

	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		panic(err)
	}

	if err := m.Down(); err != nil {
		panic(err)
	}
}
