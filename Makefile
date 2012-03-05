test:
	@go install via/via
	@via build ccache bash bash-completion
	@via install ccache bash bash-completion

all:
	via build bash ncurses pkg-config which

clean:
	@rm *.gz
	@rm *.sig
