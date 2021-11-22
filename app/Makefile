BINARY_NAME=microGoApp

build:
	@go mod vendor
	@echo "Building microGo..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "microGo built!"

run: build
	@echo "Starting microGo..."
	@./tmp/${BINARY_NAME} &
	@echo "microGo started!"

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
	@echo "Stopping microGo..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped microGo!"

restart: stop start