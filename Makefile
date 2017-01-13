GO=go
BIN=goden

all: $(BIN)

$(BIN): main.go
	$(GO) build .

install:$(BIN)
	cp $(BIN) $(GOPATH)/bin/$(BIN)

clean:
	rm $(BIN)
