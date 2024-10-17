FROM golang:1.23-alpine AS builder

COPY . /github.com/Oleg-Pro/auth
WORKDIR /github.com/Oleg-Pro/auth

RUN go mod download
RUN go build -o ./bin/auth_server cmd/auth/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Oleg-Pro/auth/bin/auth_server .

ADD .env .

CMD ["./auth_server"]