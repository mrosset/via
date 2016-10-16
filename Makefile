SRC 	= $(wildcard via//Makefile Makefile pkg/*.go via/*.go docker/Dockerfile)
BIN 	= $(GOPATH)/bin/via
CMDS	= fmt test install
REPO  = strings/via:devel

$(BIN): $(SRC)
	CGO_ENABLED=0 go install
	@git diff --quiet || echo WARNING: git tree is dirty

foo: $(BIN)

docker/via: $(BIN)
	CGO_ENABLED=0 go build -o $@

fmt:
	go fmt ./...

run:
	docker run -it strings/via:devel /bin/bash --login -o vi


dock: docker/via
	docker build -t strings/via:devel docker

clean:
	-rm docker/via
	-rm -fr root
	-rm $(BIN)

test: $(BIN)
	go test -v ./...
