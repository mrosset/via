test:
	@go install via/via
	@via build ccache
	@via install ccache

all:
	via build bash ncurses pkg-config which

clean:
	@rm *.gz
	@rm *.sig
