core:

tools:
	via build ncurses
	via build bash
	via build bzip2
	via build coreutils
	via build file
	via build findutils
	via build gawk
	via build gettext
	via build grep
	via build gzip
	via build m4
	via build make
	via build patch
	via build perl
	via build sed
	via build tar
	via build texinfo
	via build xz
	via strip_tools
	
clean:
	rm -rf cache/{builds,stages,packages}

bootstrap:
	./boostrap
