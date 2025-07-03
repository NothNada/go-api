# Etapa 1: build CGO com suporte a sqlite
FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o app .

# Etapa 2: runtime leve em Alpine
FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite
WORKDIR /
COPY --from=builder /app/app ./app
COPY data.db /data/data.db
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["./app"]
