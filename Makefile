test:
	go install via/via
	via build ccache bash
	via sign
	via install ccache bash
	#via install bash which
	#via remove bash which

all:
	via build bash ncurses pkg-config which
