# pgbr

PostgreSQL simple backup &amp; restore helper tool created for usage with Backup Repository, but can be used also standalone.

**Features:**
- `psql`, `pg_dump`, `pg_dumpall`, `pg_restore` packaged **in a single binary**, available as subcommands e.g. `pgbr psql -- -h 127.0.0.1 -d riotkit -c "SELECT 1"`
- Minimum system requirements, no extra binaries or libraries and **even no PostgreSQL client is required, just this one binary**
- Opinionated backup & restore commands basing on PostgreSQL built-in commands

**Requirements:**
- [patchelf](https://github.com/NixOS/patchelf)
- Linux


Backup
------

Selected database or all databases are dumped into a `custom formatted` file, readable by `pg_restore`.

```bash
# for single database "pbr"
pgbr db backup --password riotkit --user riotkit --db-name pbr > dump.gz

# for all databases
pgbr db backup --password riotkit --user riotkit > dump.gz
```


Restore
-------

Procedure:
1) Existing connections to selected one database, or to all databases are terminated
2) Selected database, or all databases are closed for incoming connections
3) Selected database, or all databases are recreated from backup using `pg_restore`, which uses `--clean` and `--create` by default

```bash
cat dump.gz | ./.build/pgbr db restore --password riotkit --user riotkit --connection-database=postgres
```
