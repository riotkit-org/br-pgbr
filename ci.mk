GO_VERSION=1.19
POSTGRES_VERSION=15

ci_release_snapshot:
	docker build . --build-arg POSTGRES_VERSION=${POSTGRES_VERSION} --build-arg GO_VERSION=${GO_VERSION} -t ghcr.io/riotkit-org/pgbr:latest-pg${POSTGRES_VERSION}
	docker push ghcr.io/riotkit-org/pgbr:latest-pg${POSTGRES_VERSION}

ci_release:
	docker tag ghcr.io/riotkit-org/pgbr:latest-pg${POSTGRES_VERSION} ghcr.io/riotkit-org/pgbr:$${GITHUB_REF##*/}-pg${POSTGRES_VERSION}
	docker push ghcr.io/riotkit-org/pgbr:$${GITHUB_REF##*/}-pg${POSTGRES_VERSION}

ci_rename_release_binary:
	cp .build/pgbr .build/pgbr-linux-amd64-glibc-$${GITHUB_REF##*/}-pg${POSTGRES_VERSION}
