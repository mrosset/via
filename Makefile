test: cmd
	@via build ccache

cmd: 
	@go install via/via

all:
	via build bash ncurses pkg-config which

depends:
	go get code.google.com/p/go.crypto/openpgp
	go get code.google.com/p/go.crypto/openpgp/packet

clean:
	@rm *.gz
	@rm *.sig

