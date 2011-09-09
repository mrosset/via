include $(GOROOT)/src/Make.inc

TARG=vpack
GOFILES=main.go
GOFMT=gofmt -l -w
CLEANFILES+=tmp/

include $(GOROOT)/src/Make.cmd

format:
	    ${GOFMT} .

test: all
	./${TARG} pack bash vim
	./${TARG} install bash vim
