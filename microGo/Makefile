BINARY_NAME=rapiditasApp

build:
	@go mod vendor
	@echo "Building Rapiditas..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "Rapiditas built!"

run: build
	@echo "Starting rapiditas..."
	@./tmp/${BINARY_NAME} &
	@echo "Rapiditas started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaned!"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

stop:
	@echo "Stopping rapiditas..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped Rapiditas!"

restart: stop start