.PHONY: proto server client clean install help

# Default target
help:
	@echo "📚 gRPC Demo Makefile Commands:"
	@echo "  make install    - Install required dependencies"
	@echo "  make proto      - Generate Go code from proto files"
	@echo "  make server     - Run the gRPC server"
	@echo "  make client     - Run the gRPC client"
	@echo "  make clean      - Clean generated files"
	@echo "  make all        - Generate proto and run server & client"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Dependencies installed!"

# Generate proto files
proto:
	@echo "🔨 Generating Go code from proto files..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/chat.proto
	@echo "✅ Proto files generated!"

# Run server
server:
	@echo "🚀 Starting gRPC server..."
	cd server && go run main.go

# Run client
client:
	@echo "🎯 Starting gRPC client..."
	cd client && go run main.go

# Clean generated files
clean:
	@echo "🧹 Cleaning generated files..."
	rm -f proto/*.pb.go
	@echo "✅ Clean complete!"

# Generate proto and run both server and client
all: proto
	@echo "🎬 Starting demo..."
	@echo "Run 'make server' in one terminal and 'make client' in another"
