FROM golang:1.17-alpine

# задает рабочую директорию в контейнере
WORKDIR /app

COPY ./ /app


# подгрузка зависимостей
RUN go mod download

RUN go build -o main .

ENTRYPOINT ["/app/main"]
