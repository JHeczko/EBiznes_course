# first part - compiling go server
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
COPY handlers/ handlers/
COPY models/ models/
COPY db_init/ db_init/

RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
    go build -ldflags '-extldflags "-static"' -o server.out

# second part - starting binary on light alpine version
FROM alpine:3.21

RUN apk add --no-cache sqlite && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /home/appuser/app

COPY db_start.sh db_start.sh
COPY ./db_init/data_init.sql data_init.sql
COPY --from=builder /app/server.out .

RUN chmod +x server.out && \
    chmod +x db_start.sh && \
    ./db_start.sh && \
    chmod 644 database.db && \
    chown -R appuser:appgroup /home/appuser/app

USER appuser

EXPOSE 13000
CMD ["/home/appuser/app/server.out"]