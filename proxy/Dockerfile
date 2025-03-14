# Используем официальный образ Go как базовый
FROM golang:alpine3.21 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .
# Скачиваем все зависимости
RUN go mod tidy

# Устанавливаем swag CLI для генерации документации
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Генерируем Swagger-документацию
RUN swag init -g main.go

# Собираем приложение
RUN go build -o main

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
