APP := maze
BUILD_DIR := build

.PHONY: all run install uninstall clean dvi dist tests

all:
	go build -o $(BUILD_DIR)/$(APP) ./cmd/maze

run: all
	./$(BUILD_DIR)/$(APP)

install: all

uninstall:
	rm -rf $(BUILD_DIR)

clean: uninstall
	rm -rf dist

dvi:
	@cat README.md

dist: clean
	mkdir -p dist
	tar --exclude='./build' --exclude='./dist' -czf dist/maze.tar.gz .

tests:
	go test ./...
