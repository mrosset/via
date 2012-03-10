test: cmd
	#@via build ccache
	#@via install ccache
	@via -v create http://mirrors.kernel.org/gnu/gcc/gcc-4.6.3/gcc-4.6.3.tar.gz
	#@via create http://mirrors.kernel.org/gnu/wget/wget-1.13.tar.gz
	cat ~/via/plans/wget.json

cmd: 
	@go install via/via

all:
	via build bash ncurses pkg-config which

clean:
	@rm *.gz
	@rm *.sig

#@via create http://libtorrent.rakshasa.no/downloads/libtorrent-0.13.0.tar.gz
