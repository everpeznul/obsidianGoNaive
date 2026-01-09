FROM golang:1.23-alpine AS build
WORKDIR /app

COPY app/go.mod go.sum ./
RUN go mod download

COPY .. .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/app ./cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/app /app/app
COPY --from=build /protos/configs /app/configs
COPY --from=build /gatewate/web /app/web

EXPOSE 8080
CMD ["/app/app"]
