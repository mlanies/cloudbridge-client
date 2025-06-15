.PHONY: all build clean test install build-all build-russian build-msi

BINARY_NAME=cloudbridge-client
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Build targets
PLATFORMS=linux/amd64 windows/amd64 darwin/amd64
RUSSIAN_PLATFORMS=astra/amd64 alt/amd64 rosa/amd64 redos/amd64
BUILD_DIR=build

all: clean build

build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/cloudbridge-client

build-all:
	@echo "Building for all platforms..."
	@mkdir -p ${BUILD_DIR}
	@for platform in ${PLATFORMS} ${RUSSIAN_PLATFORMS}; do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output="${BUILD_DIR}/${BINARY_NAME}-$${os}-$${arch}"; \
		if [ "$${os}" = "windows" ]; then \
			output="$${output}.exe"; \
		fi; \
		echo "Building for $${os}/$${arch}..."; \
		GOOS=$${os} GOARCH=$${arch} go build ${LDFLAGS} -o $${output} ./cmd/cloudbridge-client; \
	done

build-russian:
	@echo "Building for Russian platforms..."
	@mkdir -p ${BUILD_DIR}
	@for platform in ${RUSSIAN_PLATFORMS}; do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output="${BUILD_DIR}/${BINARY_NAME}-$${os}-$${arch}"; \
		echo "Building for $${os}/$${arch}..."; \
		GOOS=linux GOARCH=$${arch} go build ${LDFLAGS} -o $${output} ./cmd/cloudbridge-client; \
	done

clean:
	@echo "Cleaning..."
	@rm -rf bin/ ${BUILD_DIR}/

test:
	@echo "Running tests..."
	@go test -v ./...

install: build
	@echo "Installing ${BINARY_NAME}..."
	@sudo cp bin/${BINARY_NAME} /usr/local/bin/
	@sudo mkdir -p /etc/cloudbridge-client
	@sudo cp config/config.yaml /etc/cloudbridge-client/
	@sudo cp deploy/cloudbridge-client.service /etc/systemd/system/
	@sudo systemctl daemon-reload
	@sudo systemctl enable cloudbridge-client
	@sudo systemctl start cloudbridge-client

uninstall:
	@echo "Uninstalling ${BINARY_NAME}..."
	@sudo systemctl stop cloudbridge-client
	@sudo systemctl disable cloudbridge-client
	@sudo rm -f /usr/local/bin/${BINARY_NAME}
	@sudo rm -f /etc/systemd/system/cloudbridge-client.service
	@sudo systemctl daemon-reload

build-msi: build-windows
	@echo "Building MSI installer..."
	@mkdir -p ${BUILD_DIR}/msi
	@cp ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ${BUILD_DIR}/msi/${BINARY_NAME}.exe
	@cp config/config.yaml ${BUILD_DIR}/msi/
	@wixl -v -D SourceDir=${BUILD_DIR}/msi deploy/windows/cloudbridge-client.wxs -o ${BUILD_DIR}/cloudbridge-client.msi

build-windows:
	@echo "Building Windows executable..."
	@mkdir -p ${BUILD_DIR}
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ./cmd/cloudbridge-client 