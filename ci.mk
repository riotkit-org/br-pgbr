GO_VERSION=1.17
POSTGRES_VERSION=14.0

ci_check_embedded_binaries:
	./.build/pgbr pg_dump -- --version
	./.build/pgbr pg_dumpall -- --version
	./.build/pgbr pg_restore -- --version
	./.build/pgbr psql -- --version

dockerfile:
	mkdir -p .build
	cat build.Dockerfile > .build/Dockerfile
	echo "" >> .build/Dockerfile
	cat release.Dockerfile >> .build/Dockerfile

ci_release_snapshot:
	docker build . -f .build/Dockerfile --build-arg POSTGRES_VERSION=${POSTGRES_VERSION} --build-arg GO_VERSION=${GO_VERSION} -t ghcr.io/riotkit-org/pgbr:latest-pg${POSTGRES_VERSION}
	docker push ghcr.io/riotkit-org/pgbr:latest-pg${POSTGRES_VERSION}

ci_release:
	docker tag ghcr.io/riotkit-org/pgbr:latest-pg${POSTGRES_VERSION} ghcr.io/riotkit-org/pgbr:$${GITHUB_REF##*/}-pg${POSTGRES_VERSION}
	docker push ghcr.io/riotkit-org/pgbr:$${GITHUB_REF##*/}-pg${POSTGRES_VERSION}
