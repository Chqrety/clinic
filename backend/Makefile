.PHONY: run dev build test migrate-up migrate-down seed clean

# Default binary runner without Docker
run:
	@echo "🚀 Starting Clinic Accounting API..."
	go run cmd/api/main.go

dev:
	@echo "⚡ Starting Live Reload with Air..."
	air

worker:
	@echo "⚙️ Starting Background Accounting Worker..."
	go run cmd/worker/main.go

migrate-up:
	@echo "📦 Applying database migrations..."
	go run cmd/migrate/main.go -action=up

migrate-down:
	@echo "↩️ Rolling back latest migration..."
	go run cmd/migrate/main.go -action=down

seed:
	@echo "🌱 Seeding Chart of Accounts & master data..."
	go run cmd/migrate/main.go -action=seed-coa

test:
	@echo "🧪 Running all accounting unit & integration tests..."
	go test -v -race ./...

build:
	@echo "🔨 Compiling standalone production binaries..."
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/api cmd/api/main.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/worker cmd/worker/main.go
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/migrate cmd/migrate/main.go
	@echo "✅ Binaries compiled to bin/"

clean:
	@echo "🧹 Cleaning temporary files..."
	rm -rf bin/ tmp/

