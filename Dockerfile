ARG POSTGRES_VERSION
ARG GO_VERSION

# =============================================
#  Build application && embed PostgreSQL in it
# =============================================
FROM golang:${GO_VERSION} as builder

ADD . /workspace
RUN apt-get update && apt-get install make -y

RUN cd /workspace && \
    make build


# ============================================================
#  Create target image basing on a specific PostgreSQL version
# ============================================================
FROM postgres:${POSTGRES_VERSION}

RUN apt-get update && apt-get install -y gpg

COPY --from=builder /workspace/.build/pgbr /usr/bin/pgbr
RUN chmod +x /usr/bin/pgbr

ENTRYPOINT ["/usr/bin/pgbr"]
