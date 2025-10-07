#!/bin/bash
set -e

echo "Building frontend..."
cd frontend
npm install
npm run build

echo "Copying frontend build to backend..."
cd ..
rm -rf backend/frontend/dist
mkdir -p backend/frontend
cp -r frontend/dist backend/frontend/

echo "Building backend..."
cd backend
go mod download
go build -o ../tsk

echo "Build complete! Binary: ./tsk"
