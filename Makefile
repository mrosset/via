SRC 	= $(wildcard pkg/*.go via/*.go)
#BIN 	= $(GOPATH)/bin/via
BIN 	= via/via
CMDS	= fmt test install

run: docker
	docker run -t -i -v /usr/local/via:/usr/local/via -v /home/strings:/home/strings strings/via:devel bash --login -o vi

docker: root
	docker build -t strings/via:devel .

root: $(BIN)
	mkdir root
	#$(BIN) -r root install via bash coreutils
	mkdir -p root/{bin,etc,tmp}
	ln -s /usr/local/via/bin/bash root/bin/sh
	ln -s /usr/local/via/lib root/lib64
	cp -a /etc/ssl root/etc/
	cp /etc/{passwd,group} root/etc/
	tar -C root -c . | docker import - strings/via

$(BIN): $(SRC)
	make -C via
	@git diff --quiet || echo WARNING: git tree is dirty

clean:
	-rm -fr root
	-rm $(BIN)

test:
	go test -v ./...
