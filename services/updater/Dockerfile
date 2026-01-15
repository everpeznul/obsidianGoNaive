FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o updater ./services/updater/cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/updater /app/updater

EXPOSE 9002
CMD ["/app/updater"]
