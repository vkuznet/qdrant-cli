APP := qdrant-cli

.PHONY: all build clean run

all: build

build:
	go build -trimpath -o $(APP)

run:
	go run main.go

clean:
	rm -rf $(BIN)

