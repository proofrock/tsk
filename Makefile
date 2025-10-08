.PHONY: all build frontend backend clean help

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

all: build

build: frontend backend

frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build
	@echo "Copying frontend build to backend..."
	rm -rf backend/frontend/dist
	mkdir -p backend/frontend
	cp -r frontend/dist backend/frontend/

backend:
	@echo "Building backend..."
	cd backend && go mod download && go build -ldflags="-X main.Version=$(VERSION)" -o ../tsk

clean:
	@echo "Cleaning build artifacts..."
	rm -rf frontend/node_modules
	rm -rf frontend/dist
	rm -rf frontend/.vite
	rm -rf backend/frontend
	rm -f backend/tsk
	rm -f tsk
	rm -f tsk.db
	@echo "Cleanup complete!"

help:
	@echo "Available targets:"
	@echo "  make build    - Build frontend and backend (default)"
	@echo "  make frontend - Build only frontend"
	@echo "  make backend  - Build only backend"
	@echo "  make clean    - Remove all build artifacts"
	@echo "  make help     - Show this help message"
