APP_NAME=taxi-fare

build: dep
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags -installsuffix -o ${APP_NAME} github.com/mariojuzar/taxi-fare

dep:
	@echo "# Downloading Dependencies"
	@go mod download

build-run: build
	@./${APP_NAME}

run: dep
	@echo "Running taxi fare service"
	@go run github.com/mariojuzar/taxi-fare
