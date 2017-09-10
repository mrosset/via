SRC				= $(wildcard *.go Makefile pkg/*.go via/*.go docker/Dockerfile)
BIN				= $(GOPATH)/bin/via
CMDS			= fmt test install
REPO      = strings/via:devel
bash      = docker/bin/bash
btarball  = tmp/bash-4.4.tar.gz

default: $(BIN)

$(BIN): $(SRC)
	CGO_ENABLED=0 go install
	@git diff --quiet || echo WARNING: git tree is dirty

fmt:
	go fmt ./...

start:
	-docker rm -f via
	docker run --name via -it -d -e TERM=$(TERM) -v via:/home/mrosset/via -v /tmp:/tmp -v /home:/home strings/via:devel

attach: start
	docker container attach via

run:
	docker run -it -e TERM=$(TERM) -e DISPLAY=$(DISPLAY) -v /tmp:/tmp -v /tmp/.X11-unix:/tmp/.X11-unix:rw  -v /home:/home strings/via:devel /bin/bash --login -o vi

root: $(BIN)
	docker rmi -f $(REPO)
	-mkdir root
	-mkdir -p root/bin
	-ln -s /usr/local/via/bin/bash root/bin/sh
	-$(BIN) install -r root core
	-	tar -C root -c . | docker import - $(REPO)

dock: $(SRC) bash
	CGO_ENABLED=0 go build -o docker/usr/bin/via
	CGO_ENABLED=0 GOPATH=$(PWD)/docker/usr go get -v github.com/gocircuit/circuit/cmd/circuit
	docker build -t strings/via:devel docker

clean:
	-rm docker/via
	-rm -fr root
	-rm $(BIN)

test: $(BIN)
	via help
	via elf /usr/local/via/bin/bash

devel:
	bin/bdevel

bash: $(btarball) tmp/bash-4.4 tmp/bash-4.4/config.status tmp/bash-4.4/bash $(bash)

$(bash):
	mkdir -p docker/bin
	cp tmp/bash-4.4/bash $@

tmp/bash-4.4/config.status:
	cd tmp/bash-4.4; CFLAGS="-static"; ./configure

tmp/bash-4.4/bash:
	cd tmp/bash-4.4; make
	file $@

bash-clean:
	-rm tmp/bash-4.4/{bash,config.status}


$(btarball):
	mkdir -p tmp
	wget http://mirrors.kernel.org/gnu/bash/bash-4.4.tar.gz -O $(btarball)

tmp/bash-4.4:
	tar -C ./tmp -xzf $(btarball)
