# pg-simple-backup

PostgreSQL simple backup &amp; restore helper tool created for usage with Backup Repository, but can be used also standalone.

**Features:**
- `psql`, `pg_dump`, `pg_dumpall`, `pg_restore` packaged **in a single binary**, available as subcommands e.g. `pgbr psql -- -h 127.0.0.1 -d riotkit -c "SELECT 1"`
- Minimum system requirements, no extra binaries or libraries and **even no PostgreSQL client is required, just this one binary**
- Opinionated backup & restore commands basing on PostgreSQL built-in commands

**Requirements:**
- [patchelf](https://github.com/NixOS/patchelf)
