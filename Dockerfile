# Stage 1: Build the frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY clio/package.json clio/package-lock.json ./
RUN npm ci
COPY clio/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# Install build dependencies for CGO (SQLite)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the built frontend from the previous stage
COPY --from=frontend-builder /app/frontend/dist ./clio/dist
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

# Stage 3: Final image
FROM alpine:3.20.1
WORKDIR /app
COPY --from=backend-builder /app/main .
# Create directory for the database
RUN mkdir -p /data/db
EXPOSE 38080
CMD ["./main"]
