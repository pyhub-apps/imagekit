.PHONY: build build-all build-windows build-mac build-linux clean test test-coverage fmt lint run help

# 변수 설정
BINARY_NAME=imagekit
MAIN_PATH=cmd/imagekit/main.go
PACKAGE_PATH=./...

# 기본 타겟
all: build

## help: 사용 가능한 명령어 표시
help:
	@echo "사용 가능한 명령어:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: 현재 플랫폼용 바이너리 빌드
build:
	@echo "Building ${BINARY_NAME}..."
	@go build -o ${BINARY_NAME} ${MAIN_PATH}
	@echo "Build complete: ${BINARY_NAME}"

## build-all: 모든 플랫폼용 바이너리 빌드
build-all: build-windows build-mac build-linux

## build-windows: Windows용 바이너리 빌드
build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o dist/${BINARY_NAME}-windows-amd64.exe ${MAIN_PATH}
	@echo "Windows build complete"

## build-mac: macOS용 바이너리 빌드
build-mac:
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -o dist/${BINARY_NAME}-darwin-amd64 ${MAIN_PATH}
	@GOOS=darwin GOARCH=arm64 go build -o dist/${BINARY_NAME}-darwin-arm64 ${MAIN_PATH}
	@echo "macOS build complete"

## build-linux: Linux용 바이너리 빌드
build-linux:
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o dist/${BINARY_NAME}-linux-amd64 ${MAIN_PATH}
	@GOOS=linux GOARCH=arm64 go build -o dist/${BINARY_NAME}-linux-arm64 ${MAIN_PATH}
	@echo "Linux build complete"

## build-wasm: WebAssembly 빌드
build-wasm:
	@./build-wasm.sh

## clean: 빌드 결과물 삭제
clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}
	@rm -rf dist/
	@rm -f output-*.jpeg
	@rm -f output-*.png
	@echo "Clean complete"

## test: 테스트 실행
test:
	@echo "Running tests..."
	@go test -v ${PACKAGE_PATH}

## test-coverage: 커버리지와 함께 테스트 실행
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ${PACKAGE_PATH}
	@echo ""
	@echo "Detailed coverage:"
	@go test -coverprofile=coverage.out ${PACKAGE_PATH}
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## fmt: 코드 포맷팅
fmt:
	@echo "Formatting code..."
	@go fmt ${PACKAGE_PATH}
	@echo "Formatting complete"

## lint: 정적 분석 실행
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run ${PACKAGE_PATH}

## vet: go vet 실행
vet:
	@echo "Running go vet..."
	@go vet ${PACKAGE_PATH}

## mod: 모듈 의존성 정리
mod:
	@echo "Tidying modules..."
	@go mod tidy
	@echo "Modules tidied"

## run: 프로그램 실행 (도움말 표시)
run: build
	@./${BINARY_NAME} --help

## install: 시스템에 설치
install: build
	@echo "Installing ${BINARY_NAME} to /usr/local/bin..."
	@sudo cp ${BINARY_NAME} /usr/local/bin/
	@echo "Installation complete"

## uninstall: 시스템에서 제거
uninstall:
	@echo "Removing ${BINARY_NAME} from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/${BINARY_NAME}
	@echo "Uninstallation complete"

## demo: 데모 실행 (샘플 이미지 변환)
demo: build
	@echo "Running demo with sample image..."
	@./${BINARY_NAME} info sample/input-1.jpeg
	@echo ""
	@echo "Converting size to 800x600..."
	@./${BINARY_NAME} convert --width=800 --height=600 sample/input-1.jpeg demo-resized.jpeg
	@./${BINARY_NAME} info demo-resized.jpeg
	@echo ""
	@echo "Converting DPI to 72..."
	@./${BINARY_NAME} convert --dpi=72 sample/input-1.jpeg demo-dpi72.jpeg
	@./${BINARY_NAME} info demo-dpi72.jpeg
	@echo ""
	@echo "Demo complete. Check demo-*.jpeg files"

## docker-build: Docker 이미지 빌드
docker-build:
	@echo "Building Docker image..."
	@docker build -t ${BINARY_NAME}:latest .
	@echo "Docker image built: ${BINARY_NAME}:latest"

## release: 릴리즈용 빌드 (모든 플랫폼 + 압축)
release: clean build-all
	@echo "Creating release archives..."
	@mkdir -p releases
	@cd dist && tar czf ../releases/${BINARY_NAME}-windows-amd64.tar.gz ${BINARY_NAME}-windows-amd64.exe
	@cd dist && tar czf ../releases/${BINARY_NAME}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64
	@cd dist && tar czf ../releases/${BINARY_NAME}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64
	@cd dist && tar czf ../releases/${BINARY_NAME}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64
	@cd dist && tar czf ../releases/${BINARY_NAME}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64
	@echo "Release archives created in releases/"

# 기본 타겟
.DEFAULT_GOAL := help