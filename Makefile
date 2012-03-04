test:
	go install via/via
	via build which
	via install which
	#via install bash which
	#via remove bash which

all:
	via build bash ncurses pkg-config which
