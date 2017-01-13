GO=go
BIN=god

all: god

god:
	$(GO) build .

install:god
	cp $(BIN) $(GOPATH)/bin

clean:
	rm $(BIN)
