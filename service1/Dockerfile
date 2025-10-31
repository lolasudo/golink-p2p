FROM golang:1.21-alpine

WORKDIR /app

# Копируем исходный код
COPY . .

# Скачиваем зависимости и создаем go.sum
RUN go mod tidy

# Собираем приложение
RUN go build -o /service1

EXPOSE 8080

CMD ["/service1"]