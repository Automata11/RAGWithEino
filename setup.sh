#!/bin/bash

# RAG with Eino Development Setup Script

set -e

echo "🚀 Starting RAG with Eino development environment..."

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check dependencies
echo "📋 Checking dependencies..."

if ! command_exists go; then
    echo "❌ Go is not installed. Please install Go 1.19 or later."
    exit 1
fi

if ! command_exists node; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or later."
    exit 1
fi

if ! command_exists npm; then
    echo "❌ npm is not installed. Please install npm."
    exit 1
fi

echo "✅ All dependencies are installed"

# Setup backend
echo "🔧 Setting up backend..."
cd backend

if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from example..."
    cp .env.example .env
    echo "⚠️  Please edit backend/.env file with your database configurations"
fi

echo "📦 Installing backend dependencies..."
go mod download

echo "🏗️  Building backend..."
go build -o main cmd/server/main.go

cd ..

# Setup frontend
echo "🎨 Setting up frontend..."
cd frontend

echo "📦 Installing frontend dependencies..."
npm install

echo "🏗️  Building frontend..."
npm run build

cd ..

echo "✅ Setup complete!"
echo ""
echo "🔍 Next steps:"
echo "1. Configure your databases (MySQL, Redis, Milvus)"
echo "2. Update backend/.env with your database connection strings"
echo "3. Start the backend: cd backend && ./main"
echo "4. Start the frontend: cd frontend && npm run dev"
echo ""
echo "📚 For Docker deployment, run: docker-compose up -d"