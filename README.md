pgbr
====

PostgreSQL simple backup &amp; restore helper tool created for usage with Backup Repository, but can be used also standalone.

**Features:**
- Opinionated backup & restore commands basing on PostgreSQL built-in commands

**Requirements:**
- Linux (x86_64/amd64 architecture)
- PostgreSQL in desired version

Conception
----------

**Sensible defaults**

Backup & Restore should be simple and fault-tolerant, that's why this tool is automating basic things like disconnecting clients, or connecting
to database using an empty database schema during restore - we cannot restore database we connect to, also we cannot recreate from backup something that is in use.

Both `pgbr db backup` and `pgbr db restore` should work out-of-the-box with sensible defaults.

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

Passing extra arguments
-----------------------

Both `pgbr db backup` and `pgbr db restore` are supporting UNIX-like parameters passing to subprocess, which is `pg_dump`/`pg_dumpall` for **Backup** and `pg_restore` for **Restore**.

**Example**

```bash
pgbr db backup --password riotkit --user riotkit -- --role=my-role > dump.gz
```
