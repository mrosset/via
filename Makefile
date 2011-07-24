glibc-bootstrap:
	#via build_install filesystem
	#via build_install linux-api-headers

tools:
	#via build_install ncurses
	#via build_install bash
	#via build_install bzip2
	#via build_install coreutils
	#via build_install file
	#via build_install findutils
	#via build_install gawk
	#via build_install gettext
	#via build_install grep
	#via build_install gzip
	#via build_install m4
	#via build_install make
	#via build_install patch
	#via build_install perl
	#via build_install sed
	#via build_install tar
	#via build_install texinfo
	#via build_install xz
	#via build_install diffutils
	#via strip_tools
	
clean:
	rm -rf cache/builds
	rm -rf cache/stages
	rm -rf cache/packages

bootstrap:
	#./boostrap
