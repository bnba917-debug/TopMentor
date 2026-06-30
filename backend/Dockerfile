FROM golang:1.22-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /migrate ./cmd/migrate

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata wget \
    && adduser -D -u 10001 appuser

ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /server /app/server
COPY --from=builder /migrate /app/migrate
COPY migrations /app/migrations

RUN mkdir -p /data/uploads && chown -R appuser:appuser /data/uploads

USER appuser
EXPOSE 8080

CMD ["/app/server"]
