PROJ_NAME = chiasma

MAIN_PATH = cmd/main.go
BUILD_PATH = build/package/

run:
	go run $(MAIN_PATH) --follow 30s --phrase "galaxy" --output DP-3 --resolution 1920x1080 --verbose

build: clean
	go build --ldflags '-extldflags "-static"' -v -o $(BUILD_PATH)$(PROJ_NAME) $(MAIN_PATH)

clean:
	rm -rf $(BUILD_PATH)*

tests:
	go test ./...

lint:
	golangci-lint run