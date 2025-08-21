# ---------- Builder ----------
FROM golang:1.24-alpine AS builder
WORKDIR /src

# System deps
RUN apk add --no-cache git build-base

# Go module download first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy content
COPY cmd/ ./cmd/
COPY internal/ ./internal/
# COPY pkg/ ./pkg/

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/todo-server ./cmd/todo-server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/todo ./cmd/todo

# ---------- Runtime ----------
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/todo-server /app/todo-server
COPY --from=builder /out/todo /app/todo

# Expose API port
EXPOSE 8080

# Drop privileges
USER nonroot:nonroot

ENTRYPOINT ["/app/todo-server"]
