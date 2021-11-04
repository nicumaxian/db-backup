# DB Backup
`db-backup` is a CLI utility to easily create database snapshots, manage them and restore at any point for multiple configurations.

## Usage

```bash
$ db-backup --help

NAME:
   db-backup - A tool to backup and restore database easily

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.1

AUTHORS:
   Nicu Maxian <maxiannicu@gmail.com>
   Andrian Boscanean <boscanean.andrian@gmail.com>

COMMANDS:
   backup   Create a backup
   restore  Restore a backup
   config   Manage configurations
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```

## Backup Location

All backups are stored under `~/.db-backup/data/{configuration}`.

## Supported DB
![Postgres](https://cdn.iconscout.com/icon/free/png-256/postgresql-11-1175122.png)
![MySQL](https://lh3.googleusercontent.com/proxy/ISgsB-2GCUzfhYmWMUlHzJ2pZTQnDxF4Bqd8z_C0U3GBlzf62ciNdrLCagSsa82LuBbia4FR21BXQCqy6dxHcMbZcTWhKXnFOdHde9yXQU6BK2omczdhmlsn4kmS)

## Dependencies

Please make sure you have installed followings:

Postgres:
- `pg_dump`
- `psql`

MySQL:
- TBD