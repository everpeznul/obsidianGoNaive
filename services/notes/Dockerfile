FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o notes ./services/notes/cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/notes /app/notes

EXPOSE 9001
CMD ["/app/notes"]
