core:
	#via build zlib
	#via build file
	#via build zlib
	#via build file
	#via build gnupg
	#via build binutils
	#via build gmp
	#via build mpfr
	#via build mpc
	#via build gcc
	#via build sed
	#via build bzip2
	#via build pcre
	#via build glib
	#via build pkg-config
	#via build ncurses
	#via build util-linux
	#via build e2fsprogs
	#via build coreutils
	#via build iana-etc
	#via build m4
	#via build bison
	#via build procps
	#via build grep
	#via build readline
	#via build bash
	#via build libtool
	#via build gdbm
	#via build inetutils
	#via build perl
	#via build autoconf
	#via build automake
	#via build diffutils
	#via build gawk
	#via build findutils
	#via build curl
	#via build flex
	#via build gettext
	#via build groff
	#via build grub
	#via build gzip
	#via build iproute2
	#via build kbd
	#via build less
	#via build libpipeline
	#via build make
	#via build xz
	#via build man-db
	#via build module-init-tools
	#via build patch
	#via build psmisc
	#via build shadow
	#via build sysklogd
	#via build sysvinit
	#via build tar
	#via build texinfo
	#via build udev
	#via build vim

glibc-bootstrap:
	#via build filesystem
	#via build linux-api-headers

tools:
	#via build ncurses
	#via build bash
	#via build bzip2
	#via build coreutils
	#via build file
	#via build findutils
	#via build gawk
	#via build gettext
	#via build grep
	#via build gzip
	#via build m4
	#via build make
	#via build patch
	#via build perl
	#via build sed
	#via build tar
	#via build texinfo
	#via build xz
	#via build diffutils
	#via strip_tools
	
clean:
	rm -rf cache/{builds,stages,packages}
	#rm -rf /tools/*

bootstrap:
	#./boostrap
