FROM golang:1.24.2 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY src ./src

# Установка Goose
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags="-s -w" github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/app "./src/cmd/bank_api"

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /bin/app /app/app
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/src/migrations ./src/migrations

ENV PATH="/usr/local/bin:${PATH}"
RUN goose -version

CMD ["/app/app"]