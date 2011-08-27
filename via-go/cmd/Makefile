include $(GOROOT)/src/Make.inc

TARG=vpack
GOFILES=main.go
GOFMT=gofmt -l -w

include $(GOROOT)/src/Make.cmd

format:
	    ${GOFMT} .

test: all
	./${TARG} install bash vim
