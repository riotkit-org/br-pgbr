GO_VERSION=1.17
POSTGRES_VERSION=14.0

ci_build:
	docker build . -f build.Dockerfile --build-arg POSTGRES_VERSION=${POSTGRES_VERSION} --build-arg GO_VERSION=${GO_VERSION} -t build
	@docker rm -f builder 2>/dev/null
	docker create --name builder build
	docker cp builder:/workspace/.build ./
	docker cp builder:/workspace/assets ./

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
