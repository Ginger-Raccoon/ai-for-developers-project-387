# Stage 1: build frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: build backend binary
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 3: minimal runtime image
FROM alpine:3.21
WORKDIR /app
COPY --from=backend-builder /app/backend/server ./server
COPY --from=frontend-builder /app/frontend/dist ./static

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]
