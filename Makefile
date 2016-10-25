SRC 	= $(wildcard via//Makefile *.go Makefile pkg/*.go via/*.go docker/Dockerfile)
BIN 	= $(GOPATH)/bin/via
CMDS	= fmt test install
REPO  = strings/via:devel

$(BIN): $(SRC)
	CGO_ENABLED=0 go install
	@git diff --quiet || echo WARNING: git tree is dirty

foo: $(BIN)

fmt:
	go fmt ./...

run:
	docker run -it -e TERM=$(TERM) -e DISPLAY=$(DISPLAY) -v /tmp/.X11-unix:/tmp/.X11-unix:rw  -v /home:/home strings/via:devel /bin/bash --login -o vi


root: $(BIN)
	docker rmi -f $(REPO)
	-mkdir root
	-mkdir -p root/bin
	-ln -s /usr/local/via/bin/bash root/bin/sh
	-$(BIN) -r root install core
	-	tar -C root -c . | docker import - $(REPO)

dock: $(SRC)
	CGO_ENABLED=0 go build -o docker/usr/bin/via
	docker build -t strings/via:devel docker

clean:
	-rm docker/via
	-rm -fr root
	-rm $(BIN)

test: $(BIN)
	go test -v ./...
