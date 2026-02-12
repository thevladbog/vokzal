FROM golang:1.25.7-alpine AS builder

WORKDIR /app

# go-common is replaced by ../../shared/go-common; CI runs `go mod vendor` so vendor/ is in context.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o bin/fiscal cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/fiscal .
COPY config.yaml .

EXPOSE 8084

CMD ["./fiscal"]
