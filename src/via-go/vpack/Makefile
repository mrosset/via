
TARG=vpack

all: 
	go build

install:
	go install

clean:
	go clean

test: all
	./${TARG} pack bash ncdu
	./${TARG} install bash vim
