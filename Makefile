SRC 	= $(wildcard via//Makefile Makefile pkg/*.go via/*.go)
BIN 	= $(GOPATH)/bin/via
#BIN 	= via/via
CMDS	= fmt test install
REPO    = strings/via:devel

$(BIN): $(SRC)
	-rm $(BIN)
	go install ./via
	make -C via
	@git diff --quiet || echo WARNING: git tree is dirty

fmt:
	go fmt ./via/ ./pkg/
run: 
	docker run -t -i -v /etc:/etc -v /var:/var -v /tmp:/tmp -v /home:/home strings/via:devel /home/mrosset/via/root/bin/bash --login -o vi

docker: docker/Dockerfile
	docker build -t strings/via:devel docker

root: $(BIN)
	-mkdir root
	#-$(BIN) -r root install glibc bash
	-mkdir -p root/bin
	-ln -s /home/mrosset/via/root/bin/sh root/bin/sh
	#-mkdir -p root/etc root/bin root/tmp root/var/empty
	#-ln -s /home/mrosset/install/bin/pwd root/bin/pwd
	#sudo cp -a /etc/ssl root/etc/
	#sudo cp /etc/passwd root/etc/
	#sudo cp /etc/group  root/etc/
	#sudo chown -R strings root

import:
	-docker rmi -f $(REPO)
	tar -C root -c . | docker import - $(REPO)

clean:
	-rm -fr root
	-rm $(BIN)

test: $(BIN)
	$(BIN) -d build make
	#go test -v ./...
