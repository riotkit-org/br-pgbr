generate_bin_data:
	mkdir -p .build/data .build/data/bin
	cp -pr $$(whereis psql|awk '{print $$2}') .build/data/psql
	cp -pr $$(whereis pg_dumpall|awk '{print $$2}') .build/data/pg_dumpall
	cp -pr $$(whereis psql|awk '{print $$2}') .build/data/pg_dump
	cp -pr $$(whereis psql|awk '{print $$2}') .build/data/pg_restore

	./hack/get-binary-with-libs.py psql ./.build/data
	./hack/get-binary-with-libs.py pg_dump ./.build/data
	./hack/get-binary-with-libs.py pg_dumpall ./.build/data
	./hack/get-binary-with-libs.py pg_restore ./.build/data
	~/go/bin/go-bindata -o assets/main.go -pkg assets ./.build/data

clean:
	rm -rf .build/*

build:
	CGO_ENABLED=0 GO111MODULE=on go build -tags nomemcopy -o ./.build/pgbr
	chmod +x ./.build/pgbr

test_psql_in_docker:
	docker rm -f pgsb 2>/dev/null || true
	docker run -d --rm --name pgsb --entrypoint sleep debian:11-slim 99999
	docker cp ./.build/pgbr pgsb:/usr/bin/pgbr

	docker exec pgsb /bin/bash -c "apt-get update; apt-get install -f patchelf"
	docker exec pgsb /bin/bash -c "/usr/bin/pgbr"
	docker exec pgsb /bin/bash -c "find /tmp/br-pg-simple-backup"
	docker exec pgsb /bin/bash -c "export LD_LIBRARY_PATH=/tmp/br-pg-simple-backup; /tmp/br-pg-simple-backup/bin/psql --version"
	docker exec pgsb /bin/bash -c "export LD_LIBRARY_PATH=/tmp/br-pg-simple-backup; /tmp/br-pg-simple-backup/bin/pg_dump --version"
