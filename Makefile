APP := qdrant-cli
BIN := bin

.PHONY: all build clean run

all: build

build:
	@mkdir -p $(BIN)
	go build -trimpath -o $(BIN)/$(APP)

run:
	go run main.go

clean:
	rm -rf $(BIN)

