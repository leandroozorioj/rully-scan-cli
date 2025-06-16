# Rully CLI Makefile
BINARY_NAME=rully
VERSION=1.0.0

# Diretório de build
BUILD_DIR=build

# Targets de build
.PHONY: all clean build windows linux darwin

all: clean build

build: windows linux darwin

# Build para Windows
windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-386.exe main.go

# Build para Linux
linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-386 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go

# Build para macOS
darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go

# Build para desenvolvimento local
dev:
	@echo "Building for development..."
	go build -o $(BINARY_NAME) main.go

# Executar testes
test:
	go test -v ./...

# Instalar dependências
deps:
	go mod download
	go mod tidy

# Limpar builds
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)

# Executar o programa
run:
	go run main.go

# Mostrar help
help:
	@echo "Rully CLI Makefile"
	@echo ""
	@echo "Targets disponíveis:"
	@echo "  all      - Build para todas as plataformas"
	@echo "  build    - Build para todas as plataformas"
	@echo "  windows  - Build apenas para Windows"
	@echo "  linux    - Build apenas para Linux"
	@echo "  darwin   - Build apenas para macOS"
	@echo "  dev      - Build para desenvolvimento local"
	@echo "  deps     - Instalar dependências"
	@echo "  test     - Executar testes"
	@echo "  clean    - Limpar builds"
	@echo "  run      - Executar o programa"
	@echo "  help     - Mostrar esta ajuda"