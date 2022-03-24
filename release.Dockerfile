# Intentionally we do not use Alpine Linux there, to ensure binaries compatibility with at least Debian-like systems. We test on CI on Ubuntu and produce Debian images.

# =====================
#  Create target image
# =====================
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y patchelf

COPY --from=builder /workspace/.build/pgbr /usr/bin/pgbr
ENTRYPOINT ["/usr/bin/pgbr"]
