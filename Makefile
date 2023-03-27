

all: install test

install:
	go install ./tex2mnemosyne/...

test:
	go test -v ./...
