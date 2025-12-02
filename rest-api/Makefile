APP_NAME=REST_API
BIN_NAME=gobin
BUILD_DIR=./bin
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

run:
	@echo "Running the server"
	@go run main.go || true

deps:
	@echo "Installing dependancies"
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./... 

build:
	@echo "Building the project"
	@mkdir -p ${BUILD_DIR}
	@go build -o $(BUILD_DIR)/$(bin) main.go
	@echo "Build complete: $(BUILD_DIR)"

clean:
	@echo "Cleaning up"
	@rm -rf ${BUILD_DIR}
	@echo "Clean complete"

stop:
	@echo "Stopping server..."
	@pkill -f "go run main.go" || echo "No process found"

migrate-up:
	@echo "Running database migration"
	@dbmate up
	@echo "Migration applied"

migrate-down:
	@echo "Rolling back last migration"
	@dbmate down
	@echo "Migration rolled back"

help:
	@echo "Available commands"
	@echo "make run		- Run the server"
	@echo "make deps	- Install Dependancies"
	@echo "make fmt		- Formats the code"
	@echo "make build	- Build the binary"
	@echo "make clean	- Clean up the binary"
	@echo "make stop 	- stop the running server/s"
