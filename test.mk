
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

test_dump:
	./.build/pgbr db backup -P riotkit -U riotkit -d pbr
