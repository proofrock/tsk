# Build frontend
FROM node:latest AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:latest AS backend-builder

ARG VERSION=dev

RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=1 go build -ldflags="-X main.Version=${VERSION}" -o tsk

# Runtime
FROM debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates libsqlite3-0 && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=backend-builder /app/tsk .

EXPOSE 8080

CMD ["./tsk"]
