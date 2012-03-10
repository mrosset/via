cmd: 
	@go install via/via

test: cmd
	#@via build ccache
	#@via install ccache
	@via create http://mirrors.kernel.org/gnu/wget/wget-1.13.tar.gz

all:
	via build bash ncurses pkg-config which

clean:
	@rm *.gz
	@rm *.sig

#@via create http://libtorrent.rakshasa.no/downloads/libtorrent-0.13.0.tar.gz
