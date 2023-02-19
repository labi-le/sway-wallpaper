PROJ_NAME = hw

MAIN_PATH = cmd/main.go
BUILD_PATH = build/package/

INSTALL_PATH = /usr/bin/

run:
	go run $(MAIN_PATH) -follow 5s -search-phrase galaxy

build: clean
	go build --ldflags '-extldflags "-static"' -v -o $(BUILD_PATH)$(PROJ_NAME) $(MAIN_PATH)

release: clean
	goreleaser --rm-dist

install:
	make build
	sudo cp $(BUILD_PATH)$(PROJ_NAME) $(INSTALL_PATH)$(PROJ_NAME)

uninstall:
	sudo rm $(INSTALL_PATH)$(PROJ_NAME)

clean:
	rm -rf $(BUILD_PATH)*

tests:
	go test ./...

lint:
	golangci-lint run