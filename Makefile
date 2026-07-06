BINARY := foundry-plans
EXEC   := $(HOME)/.local/bin
INSTALL := $(EXEC)/$(BINARY)

.PHONY: build install run clean test

build:
	go build -o $(BINARY)

install: build
	@mkdir -p $(EXEC)
	cp $(BINARY) $(INSTALL)
	@echo "installed $(INSTALL)"

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test ./...
