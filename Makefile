include test.mk
include ci.mk

copy_libs_and_executables:
	mkdir -p .build/data .build/data/bin
	cp -p $$(whereis psql|awk '{print $$2}') .build/data/psql
	cp -p $$(whereis pg_dumpall|awk '{print $$2}') .build/data/pg_dumpall
	cp -p $$(whereis pg_dump|awk '{print $$2}') .build/data/pg_dump
	cp -p $$(whereis pg_restore|awk '{print $$2}') .build/data/pg_restore

	./hack/get-binary-with-libs.py psql ./.build/data
	./hack/get-binary-with-libs.py pg_dump ./.build/data
	./hack/get-binary-with-libs.py pg_dumpall ./.build/data
	./hack/get-binary-with-libs.py pg_restore ./.build/data

generate_bin_data:
	~/go/bin/go-bindata -o assets/main.go -pkg assets ./.build/data

clean:
	rm -rf .build/*

assets_build: ## Builds PostgreSQL assets using docker
	docker build . -f build.Dockerfile --build-arg POSTGRES_VERSION=${POSTGRES_VERSION} --build-arg GO_VERSION=${GO_VERSION} -t build
	@docker rm -f builder 2>/dev/null
	docker create --name builder build
	docker cp builder:/workspace/.build ./
	docker cp builder:/workspace/assets ./

build:
	CGO_ENABLED=0 GO111MODULE=on go build -tags nomemcopy -o ./.build/pgbr
	chmod +x ./.build/pgbr

test:
	mkdir -p .build/data/bin
	cp -p .build/data/pg_* .build/data/bin/
	cp -p .build/data/psql .build/data/bin/psql
	go test -v ./...
