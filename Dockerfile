FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
ENV NODE_ENV=production
COPY clio/package.json clio/package-lock.json ./
RUN npm install --include=dev
COPY clio/ ./
RUN npm run build

FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY --from=frontend-builder /app/frontend/dist ./clio/dist

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:3.20.1
WORKDIR /app
COPY --from=backend-builder /app/main .
COPY --from=backend-builder /app/internal/database/migrations/ /app/internal/database/migrations/
COPY --from=backend-builder /app/clio/dist /app/clio/dist
RUN mkdir -p /data/db

ENV APP_ENV=production
ENV APP_PORT=38080
EXPOSE 38080
CMD ["./main"]
