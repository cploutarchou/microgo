BINARY_NAME=MicroGOApp

build:
	@echo "Building MicroGO..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "MicroGO built!"

run: build
	@echo "Starting MicroGO..."
	@./tmp/${BINARY_NAME} &
	@echo "MicroGO started!"

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
	@echo "Stopping MicroGO..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped MicroGO!"

restart: stop start
