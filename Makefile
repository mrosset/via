SRC 	= $(wildcard via//Makefile Makefile pkg/*.go via/*.go)
BIN 	= $(GOPATH)/bin/via
CMDS	= fmt test install
REPO  = strings/via:devel

$(BIN): $(SRC)
	-rm $(BIN)
	CGO_ENABLED=0 go install
#make -C via
	@git diff --quiet || echo WARNING: git tree is dirty

fmt:
	go fmt ./via/ ./pkg/

run:
	docker run -it -v /var:/var -v /tmp:/tmp -v /home:/home strings/via:devel /usr/local/via/bin/bash --login -o vi

dock:
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
	-rm -fr root
	-rm $(BIN)

test: $(BIN)
	$(BIN) -d build ccache
#go test -v ./...
