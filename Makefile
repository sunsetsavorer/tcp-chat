APP_NAME="tcp-chat"
MAIN_PATH="./cmd/main.go"

run:
	go run ${MAIN_PATH}

build:
	@echo "🔧 Building app..."
	go build -o ./build/${APP_NAME} ${MAIN_PATH}
	@echo "✅ Done!"

clean:
	rm -rf build
	@echo "✅ Done!"

help:
	@echo "Available commands:"
	@echo "  make build   - собрать бинарник"
	@echo "  make run     - запустить приложение"
	@echo "  make clean   - удалить бинарники"