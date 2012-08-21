SRC 	= $(wildcard pkg/*.go via/*.go)
BIN 	= $(GOPATH)/bin/via
CMDS	= fmt test install

$(BIN): $(SRC)
	go fmt ./...
	go test ./...
	go install ./...
	@git diff --quiet || echo WARNING: git tree is dirty

clean:
	-rm $(BIN)
