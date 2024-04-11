# Используйте официальный образ Golang как базовый
FROM golang:1.22.2

# Установите рабочий каталог в контейнере
WORKDIR /app

# Копируйте файлы go.mod и go.sum в контейнер
COPY go.mod go.sum ./

# Скачайте зависимости
RUN go mod download

# Копируйте исходный код и конфигурационные файлы в контейнер
COPY . .

# Сборка приложения
RUN go build -o main ./cmd/sso/main.go

# Запуск приложения с параметрами командной строки
CMD ["go", "run", "./cmd/sso/main.go", "--config=./internal/config/config_local.yaml"]
