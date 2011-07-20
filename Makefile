test:
	./via pbuild binutils-bootstrap
	./via pbuild gcc-bootstrap
	./via pbuild linux-api-bootstrap
	./via pbuild glibc-bootstrap
	#./via adjust_gcc
	#./via pbuild ncurses-bootstrap
	#./via pbuild ncdu-bootstrap
	#./via pbuild busybox-bootstrap
