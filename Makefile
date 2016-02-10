SRC 	= $(wildcard via//Makefile Makefile pkg/*.go via/*.go)
BIN 	= $(GOPATH)/bin/via
#BIN 	= via/via
CMDS	= fmt test install
REPO    = strings/via:devel

$(BIN): $(SRC)
	-rm $(BIN)
	make -C via
	@git diff --quiet || echo WARNING: git tree is dirty

run: 
	docker run -t -i -v /home/strings:/home/strings strings/via:devel bash --login -o vi

docker:
	docker build -t strings/via:devel .

root: $(BIN)
	-mkdir root
	-$(BIN) -r root install devel
	mkdir -p root/{etc,tmp}
	cp -a /etc/ssl root/etc/
	cp /etc/{passwd,group} root/etc/
	ldconfig -r root/

import:
	-docker rmi -f $(REPO)
	tar -C root -c . | docker import - $(REPO)

clean:
	-rm -fr root
	-rm $(BIN)

test: $(BIN)
	$(BIN) -d build ccache
	#go test -v ./...
