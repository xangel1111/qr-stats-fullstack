# ---- build stage ----
FROM golang:1.26-alpine AS build

WORKDIR /app

# Cache dependencies first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# ---- run stage ----
FROM alpine:3.20

# Run as a non-root user.
RUN adduser -D -H appuser

COPY --from=build /server /server

USER appuser
EXPOSE 8080

ENTRYPOINT ["/server"]
