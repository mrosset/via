
TARG=vpack

all: 
	go build

install:
	go install

clean:
	go clean

test: all
	./${TARG} pack bash vim
	./${TARG} install bash vim
