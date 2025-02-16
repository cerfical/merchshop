FROM golang:alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/api/main.go main.go
COPY internal internal
RUN go build -tags viper_bind_struct -v -o bin/app main.go

FROM alpine:latest

COPY --from=build /app/bin/ /usr/local/bin/
COPY internal/lib/postgres/migrations migrations

ENV MERCHSHOP_API_SERVER_HOST=0.0.0.0 \
    MERCHSHOP_API_SERVER_PORT=80

ENTRYPOINT ["app"]
