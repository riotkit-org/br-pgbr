ARG POSTGRES_VERSION
ARG GO_VERSION

# ========================================
#  Collect specific version of PostgreSQL
# ========================================
FROM postgres:${POSTGRES_VERSION}-alpine as postgres

RUN apk update && apk add python3 make bash util-linux

ADD . /workspace
WORKDIR /workspace

# collect binaries and libraries in .build/data directory
RUN cd /workspace && make copy_libs_and_executables


# =============================================
#  Build application && embed PostgreSQL in it
# =============================================
FROM golang:${GO_VERSION}-alpine as builder

ADD . /workspace
RUN go install -a -v github.com/go-bindata/go-bindata/...@latest

COPY --from=postgres /workspace/.build /workspace/.build

RUN cd /workspace \
    make generate_bin_data build
