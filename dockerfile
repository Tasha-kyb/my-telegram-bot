FROM golang:1.24

WORKDIR /app

# копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# копируем все файлы проекта
COPY . .

# собираем приложение
RUN go build -o main ./cmd

# указываем команду запуска
CMD ["./main"]