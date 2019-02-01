SRC	  = $(wildcard *.go Makefile pkg/plugin/*.go pkg/*.go via/*.go docker/Dockerfile)
PLG       = $(wildcard plugins/*.go)
BIN	  = $(GOPATH)/bin/via
CMDS	  = fmt test install
REPO      = strings/via:devel
bash      = docker/bin/bash
btarball  = tmp/bash-4.4.tar.gz

export CGO_ENABLED=1
export PREFIX=/opt/via

default: $(BIN)

devel: default
	rm -rf /opt/via/*; rm -rf ~/src/via/publish
	$(BIN) install -y devel
	echo
	du -hs publish
	du -hs $(PREFIX)/

run: default
	$(BIN) help

$(BIN): $(SRC)
	go build -o $(BIN)
	@git diff --quiet || echo WARNING: git tree is dirty
	strip $(BIN)

fmt:
	go fmt ./...

start:
	-docker rm -f via
	docker run --privileged --name via -it -d -e DISPLAY=$(DISPLAY) -e TERM=$(TERM) -v /tmp:/tmp -v /home:/home strings/via:devel

attach: start
	docker container attach via

dock: $(SRC) bash
	go build -v -ldflags "-linkmode external -extldflags -static" -o docker/usr/bin/via
	docker build -t strings/via:devel docker

clean:
	@-rm docker/via
	-rm -fr ./root
	-rm $(BIN)
	-rm -rf ./tmp/bash-4.4

rebuild: clean default

test:
	go test -v ./...
	#go test -run TestPlugin* -v ./pkg/...

.NOTPARALLEL:

bash: $(btarball) tmp/bash-4.4 tmp/bash-4.4/config.status tmp/bash-4.4/bash $(bash)

$(bash):
	mkdir -p docker/bin
	cp -p tmp/bash-4.4/bash $@

tmp/bash-4.4/config.status:
	cd tmp/bash-4.4; CFLAGS="-static"; ./configure --enable-static-link -q

tmp/bash-4.4/bash:
	$(MAKE) -C tmp/bash-4.4

bash-clean:
	-rm $(bash) tmp/bash-4.4/{bash,config.status}


$(btarball):
	mkdir -p tmp
	wget http://mirrors.kernel.org/gnu/bash/bash-4.4.tar.gz -O $(btarball)

tmp/bash-4.4:
	tar -C ./tmp -xzf $(btarball)
