include test.mk

generate_bin_data:
	mkdir -p .build/data .build/data/bin
	cp -pr $$(whereis psql|awk '{print $$2}') .build/data/psql
	cp -pr $$(whereis pg_dumpall|awk '{print $$2}') .build/data/pg_dumpall
	cp -pr $$(whereis pg_dump|awk '{print $$2}') .build/data/pg_dump
	cp -pr $$(whereis pg_restore|awk '{print $$2}') .build/data/pg_restore

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

