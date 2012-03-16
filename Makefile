test: cmd
	@via build ccache

cmd: 
	@go install via/via

all:
	via build bash ncurses pkg-config which

clean:
	@rm *.gz
	@rm *.sig
