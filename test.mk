
test_psql_in_docker:
	docker rm -f pgsb 2>/dev/null || true
	docker run -d --rm --name pgsb --entrypoint sleep debian:11-slim 99999
	docker cp ./.build/pgbr pgsb:/usr/bin/pgbr

	docker exec pgsb /bin/bash -c "apt-get update; apt-get install -f patchelf"
	docker exec pgsb /bin/bash -c "/usr/bin/pgbr"
	docker exec pgsb /bin/bash -c "find /tmp/br-pg-simple-backup"
	docker exec pgsb /bin/bash -c "export LD_LIBRARY_PATH=/tmp/br-pg-simple-backup; /tmp/br-pg-simple-backup/bin/psql --version"
	docker exec pgsb /bin/bash -c "export LD_LIBRARY_PATH=/tmp/br-pg-simple-backup; /tmp/br-pg-simple-backup/bin/pg_dump --version"

test_postgres_container:
	@docker rm -f br-postgres || true
	docker run -d --rm --name br-postgres -p 5432:5432 -e POSTGRES_USER=riotkit -e POSTGRES_PASSWORD=riotkit -e POSTGRES_DB=pbr postgres:13.6
	while ! pg_isready -h 127.0.0.1; do sleep 0.5; done

test_import_data: ## Imports bigger amount of simple structured data + real application tables with almost empty data
	# dataset that is very simple structured, but adds in a loop a lot of records
	cat ./hack/sql/create-test-data.sql | docker exec -i -e PGPASSWORD=riotkit br-postgres psql -U riotkit pbr

	# dataset of a real application, a skeleton of tables but almost no data
	curl -s https://raw.githubusercontent.com/Platformus/Platformus-Sample-Ecommerce/681b898d7822bbc3ed8dd4d9a1c1176fddfead7b/postgresql.sql --output ./.build/example-structured-data.sql
	cat ./.build/example-structured-data.sql | docker exec -i -e PGPASSWORD=riotkit br-postgres psql -U riotkit pbr

test_dump:
	mkdir -p .build/.test-backups
	@./.build/pgbr db backup -P riotkit -U riotkit -d pbr > .build/.test-backups/dump.gz

test_restore:
	@cat .build/.test-backups/dump.gz | ./.build/pgbr db restore -P riotkit -U riotkit


# todo: pg_restore: error: could not execute query: ERROR:  database "pbr" is being accessed by other users
