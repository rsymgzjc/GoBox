FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev

COPY go.mod ./
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o gobox ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/gobox ./gobox
COPY .env.example ./.env
EXPOSE 8080
CMD ["./gobox"]
