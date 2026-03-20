FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -o /shortener ./cmd/main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /shortener .
COPY --from=builder /app/web ./web
EXPOSE 8080
ENTRYPOINT ["./shortener"]
