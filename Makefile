SRC 	= $(wildcard via//Makefile Makefile pkg/*.go via/*.go)
BIN 	= $(GOPATH)/bin/via
CMDS	= fmt test install
REPO  = strings/via:devel

$(BIN): $(SRC)
	-rm $(BIN)
	CGO_ENABLED=0 go install
#make -C via
	@git diff --quiet || echo WARNING: git tree is dirty


docker/via:
	CGO_ENABLED=0 go build -o $@

fmt:
	go fmt ./via/ ./pkg/

run:
	docker run -it strings/via:devel /bin/bash --login -o vi


dock: clean docker/via
	docker build -t strings/via:devel docker

root: $(BIN)
	-mkdir root
	-$(BIN) -r root install devel
	-mkdir -p root/bin root/etc
	-cp etc/* root/etc
	-ln -s /usr/local/via/bin/sh root/bin/sh
	tar -C root -c . | docker import - $(REPO)
	#docker rmi -f $(REPO)

clean:
	-rm docker/via
	-rm -fr root
	-rm $(BIN)

test: $(BIN)
	$(BIN) -d build ccache
#go test -v ./...
