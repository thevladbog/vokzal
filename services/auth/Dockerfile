FROM golang:1.25.7-alpine AS builder

WORKDIR /build

# go-common is replaced by ../../shared/go-common; CI runs `go mod vendor` so vendor/ is in context.
COPY . .

# Собрать бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o auth-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Скопировать бинарник
COPY --from=builder /build/auth-service .
COPY --from=builder /build/config.yaml .

# Создать пользователя
RUN addgroup -g 1000 vokzal && \
    adduser -D -u 1000 -G vokzal vokzal && \
    chown -R vokzal:vokzal /app

USER vokzal

EXPOSE 8081

CMD ["./auth-service"]
