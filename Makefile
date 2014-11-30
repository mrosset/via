SRC 	= root $(wildcard pkg/*.go via/*.go)
#BIN 	= $(GOPATH)/bin/via
BIN 	= via/via
CMDS	= fmt test install

docker: root
	docker build .

root: $(BIN)
	mkdir root
	$(BIN) -r root install glibc dash coreutils
	mkdir -p root/bin
	ln -s /usr/local/via/bin/dash root/bin/sh
	ln -s /usr/local/via/lib root/lib64
	mkdir root/etc
	mkdir root/tmp
	cp -a /etc/ssl root/etc/
	cp /etc/passwd root/etc/
	mkdir -p root/root/via
	git clone $(HOME)/via/plans root/root/via/plans
	tar -C root -c . | docker import - strings/via

$(BIN): $(SRC)
	make -C via
	@git diff --quiet || echo WARNING: git tree is dirty

clean:
	-rm -fr root
	-rm $(BIN)
