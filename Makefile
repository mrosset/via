SRC				= $(wildcard *.go Makefile pkg/*.go via/*.go docker/Dockerfile)
BIN				= $(GOPATH)/bin/via
CMDS			= fmt test install
REPO      = strings/via:devel
bash      = docker/bin/bash
btarball  = tmp/bash-4.4.tar.gz

default: $(BIN)

$(BIN): $(SRC)
	CGO_ENABLED=0 go build -o $(BIN)
	@git diff --quiet || echo WARNING: git tree is dirty
	strip $(BIN)
	file $(BIN)

fmt:
	go fmt ./...

start:
	-docker rm -f via
	docker run --name via -it -d -e TERM=$(TERM) -v via:/via -v /tmp:/tmp -v /home:/home strings/via:devel

attach: start
	docker container attach via

run:
	-docker rm bash
	docker run --name bash -it -e TERM=$(TERM) -e DISPLAY=$(DISPLAY) -v via:/via -v /tmp:/tmp -v /tmp/.X11-unix:/tmp/.X11-unix:rw  -v /home:/home strings/via:devel /bin/bash --login -o vi

dock: $(SRC) bash
	CGO_ENABLED=0 go build -o docker/usr/bin/via
	docker build -t strings/via:devel docker

clean:
	-rm docker/via
	-rm -fr root
	-rm $(BIN)

rebuild: clean default

test:
	go test ./pkg/...

bash: $(btarball) tmp/bash-4.4 tmp/bash-4.4/config.status tmp/bash-4.4/bash $(bash)

$(bash):
	mkdir -p docker/bin
	cp -p tmp/bash-4.4/bash $@

tmp/bash-4.4/config.status:
	cd tmp/bash-4.4; CFLAGS="-static"; ./configure -q

tmp/bash-4.4/bash:
	$(MAKE) -C tmp/bash-4.4

bash-clean:
	-rm $(bash) tmp/bash-4.4/{bash,config.status}


$(btarball):
	mkdir -p tmp
	wget http://mirrors.kernel.org/gnu/bash/bash-4.4.tar.gz -O $(btarball)

tmp/bash-4.4:
	tar -C ./tmp -xzf $(btarball)
