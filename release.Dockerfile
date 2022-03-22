# =====================
#  Create target image
# =====================
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y patchelf

COPY --from=builder /workspace/.build/pgbr /usr/bin/pgbr
ENTRYPOINT ["/usr/bin/pgbr"]
