# Go 1.24 to match your go.mod
FROM golang:1.24-alpine

WORKDIR /app
RUN apk add --no-cache git

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build the ROOT package (where main.go lives)
COPY . .
RUN CGO_ENABLED=0 go build -o /app/todo-app . \
 && chmod +x /app/todo-app

EXPOSE 8080
# Run the CLI; we'll pass subcommands via docker-compose
ENTRYPOINT ["/app/todo-app"]
