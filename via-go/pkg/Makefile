include $(GOROOT)/src/Make.inc

TARG=via
GOFILES=via.go plan.go archive.go manifest.go repo.go json.go http.go
GOFMT=gofmt -l -w
CLEANFILES+=*.bz2 tmp/
include $(GOROOT)/src/Make.pkg

format:
	    ${GOFMT} .
