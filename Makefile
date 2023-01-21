include test.mk
include ci.mk

clean:
	rm -rf .build/*

build:
	CGO_ENABLED=0 GO111MODULE=on go build -tags nomemcopy -o ./.build/pgbr
	chmod +x ./.build/pgbr

test:
	export PGBR_USE_CONTAINER=true; \
	export PGBR_CONTAINER_IMAGE=bitnami/postgresql; \
	export POSTGRES_VERSION=${POSTGRES_VERSION}; \
	go test -v ./...
