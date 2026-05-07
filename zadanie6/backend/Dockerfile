# first part - compiling go server
FROM golang:1.26-alpine AS builder
# Instalujemy narzędzia do kompilacji C (potrzebne dla CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# KLUCZOWE ZMIANY:
# 1. CGO_ENABLED=1 (bo standardowy sqlite go potrzebuje)
# 2. GOARCH=arm64 (bo taką masz architekturę wg. błędu)
# 3. -ldflags '-extldflags "-static"' (wymuszamy wrzucenie wszystkiego do środka)
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
    go build -ldflags '-extldflags "-static"' -o server.out

# second part - starting binary on light alpine version
FROM alpine:latest

RUN apk add --no-cache sqlite

WORKDIR /root/app

COPY db_start.sh db_start.sh
COPY ./db_init/data_init.sql data_init.sql
COPY --from=builder /app/server.out .

RUN chmod +x server.out && \
    chmod +x db_start.sh && \
    ./db_start.sh && \
    chmod 777 database.db

EXPOSE 13000
CMD ["/root/app/server.out"]