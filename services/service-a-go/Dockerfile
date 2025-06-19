# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.24.3 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o service-a .

# --- Runtime stage ---
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/service-a /service-a
EXPOSE 8081 50051
ENTRYPOINT ["/service-a"]
