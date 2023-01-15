include test.mk
include ci.mk

PG_DATA=assets/.build/data

copy_libs_and_executables:
	mkdir -p ${PG_DATA}
	cp -p /usr/lib/postgresql/${POSTGRES_VERSION}/bin/psql ${PG_DATA}/psql
	cp -p /usr/lib/postgresql/${POSTGRES_VERSION}/bin/pg_dumpall ${PG_DATA}/pg_dumpall
	cp -p /usr/lib/postgresql/${POSTGRES_VERSION}/bin/pg_dump ${PG_DATA}/pg_dump
	cp -p /usr/lib/postgresql/${POSTGRES_VERSION}/bin/pg_restore ${PG_DATA}/pg_restore

	./hack/get-binary-with-libs.py /usr/lib/postgresql/${POSTGRES_VERSION}/bin/psql ./${PG_DATA}
	./hack/get-binary-with-libs.py /usr/lib/postgresql/${POSTGRES_VERSION}/bin/pg_dump ./${PG_DATA}
	./hack/get-binary-with-libs.py /usr/lib/postgresql/${POSTGRES_VERSION}/bin/pg_dumpall ./${PG_DATA}
	./hack/get-binary-with-libs.py /usr/lib/postgresql/${POSTGRES_VERSION}/bin/pg_restore ./${PG_DATA}

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

test_runs:
	if [[ $$CI == "true" ]]; then \
  		sudo /bin/bash -c 'echo "pgsqluser:x:$$(id -u):$$(id -g)::/home:/bin/bash" >> /etc/passwd'; \
  	fi

	./.build/pgbr pg_dump -- --help
	./.build/pgbr pg_dumpall -- --help
	./.build/pgbr psql -- --help

test: test_runs
	go test -v ./...
