FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o gateway ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/gateway /app/gateway
COPY --from=build /app/gatewate/web /app/web

EXPOSE 8080
CMD ["/app/gateway"]
