FROM golang:1.23.4-alpine AS builder

COPY . /github.com/mukhinfa/auth/source

WORKDIR /github.com/mukhinfa/auth/source

RUN go mod download
RUN go build -o ./bin/auth_service cmd/main.go

FROM alpine:latest

WORKDIR /root
COPY --from=builder /github.com/mukhinfa/auth/source/bin/auth_service .

CMD ["./auth_service"]
