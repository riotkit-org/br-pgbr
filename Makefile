include test.mk
include ci.mk

clean:
	rm -rf .build/*

build:
	CGO_ENABLED=0 GO111MODULE=on go build -tags nomemcopy -o ./.build/pgbr
	chmod +x ./.build/pgbr

test:
	go test -v ./...
