pgbr
====

PostgreSQL simple backup &amp; restore helper tool created for usage with Backup Repository, but can be used also standalone.

**Features:**
- `psql`, `pg_dump`, `pg_dumpall`, `pg_restore` packaged **in a single binary**, available as subcommands e.g. `pgbr psql -- -h 127.0.0.1 -d riotkit -c "SELECT 1"`
- Minimum system requirements, no extra binaries or libraries and **even no PostgreSQL client is required, just this one binary**
- Opinionated backup & restore commands basing on PostgreSQL built-in commands

**Requirements:**
- [patchelf](https://github.com/NixOS/patchelf)
- Linux

Conception
----------

**Single-binary**

`pgbr` binary has compiled PostgreSQL tools, including libc (musl/glibc) and dynamic libraries, everything is unpacked into temporary directory, then `patchelf` patches
interpreter to match `ld-linux` or `ld-musl` at selected path and `pgbr` is invoking a process with extra `LD_LIBRARY_PATH` environment set.

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

psql & pg_dump & pg_dumpall & pg_restore
----------------------------------------

- psql: `pgbr psql -- ...`
- pg_dump: `pgbr pg_dump -- ...`
- pg_dumpall: `pgbr pg_dumpall -- ...`
- pg_restore: `pgbr pg_restore -- ...`

```bash
pgbr psql -- --help

#psql is the PostgreSQL interactive terminal.
# 
#Usage:
#  psql [OPTION]... [DBNAME [USERNAME]]
# 
#General options:
#  -c, --command=COMMAND    run only single command (SQL or internal) and exit
#  -d, --dbname=DBNAME      database name to connect to (default: "damian")
#  -f, --file=FILENAME      execute commands from file, then exit
#  -l, --list               list available databases, then exit
#  -v, --set=, --variable=NAME=VALUE
#                           set psql variable NAME to VALUE
#                           (e.g., -v ON_ERROR_STOP=1)
#  -V, --version            output version information, then exit
#  -X, --no-psqlrc          do not read startup file (~/.psqlrc)
#  -1 ("one"), --single-transaction
#                           execute as a single transaction (if non-interactive)
#  -?, --help[=options]     show this help, then exit
#      --help=commands      list backslash commands, then exit
#      --help=variables     list special variables, then exit
#
#Input and output options:
#  -a, --echo-all           echo all input from script
#  -b, --echo-errors        echo failed commands
#  -e, --echo-queries       echo commands sent to server
#  -E, --echo-hidden        display queries that internal commands generate
#  -L, --log-file=FILENAME  send session log to file
#  -n, --no-readline        disable enhanced command line editing (readline)
#  -o, --output=FILENAME    send query results to file (or |pipe)
#  -q, --quiet              run quietly (no messages, only query output)
#  -s, --single-step        single-step mode (confirm each query)
#  -S, --single-line        single-line mode (end of line terminates SQL command)
#
#Output format options:
#  -A, --no-align           unaligned table output mode
#      --csv                CSV (Comma-Separated Values) table output mode
#  -F, --field-separator=STRING
#                           field separator for unaligned output (default: "|")
#  -H, --html               HTML table output mode
#  -P, --pset=VAR[=ARG]     set printing option VAR to ARG (see \pset command)
#  -R, --record-separator=STRING
#                           record separator for unaligned output (default: newline)
#  -t, --tuples-only        print rows only
#  -T, --table-attr=TEXT    set HTML table tag attributes (e.g., width, border)
#  -x, --expanded           turn on expanded table output
#  -z, --field-separator-zero
#                           set field separator for unaligned output to zero byte
#  -0, --record-separator-zero
#                           set record separator for unaligned output to zero byte
#
#Connection options:
#  -h, --host=HOSTNAME      database server host or socket directory (default: "local socket")
#  -p, --port=PORT          database server port (default: "5432")
#  -U, --username=USERNAME  database user name (default: "damian")
#  -w, --no-password        never prompt for password
#  -W, --password           force password prompt (should happen automatically)
#
#For more information, type "\?" (for internal commands) or "\help" (for SQL
#commands) from within psql, or consult the psql section in the PostgreSQL
#documentation.
#
#Report bugs to <pgsql-bugs@lists.postgresql.org>.
#PostgreSQL home page: <https://www.postgresql.org/>
```

Passing extra arguments
-----------------------

Both `pgbr db backup` and `pgbr db restore` are supporting UNIX-like parameters passing to subprocess, which is `pg_dump`/`pg_dumpall` for **Backup** and `pg_restore` for **Restore**.

**Example**

```bash
pgbr db backup --password riotkit --user riotkit -- --role=my-role > dump.gz
```
