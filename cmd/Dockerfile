FROM golang:1.23.1-alpine3.20 AS build

RUN apk add --no-cache git

WORKDIR /var/backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine:edge as prod

WORKDIR /app
COPY --from=build /var/backend/main .
COPY --from=build /var/backend/.env .
COPY --from=build /var/backend/configs ./configs

EXPOSE 8080

ENTRYPOINT ./main