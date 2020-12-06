# backup

backup is a simple and configurable wrapper of rsync to easily perform backups
and restore backups.

**Disclaimer**: this tool is not a production ready tool and you should use with
care and knowledge ! Always use `-f` flag before an operation to make a dry run.

## Installation

```
go get -u github.com/nikkolasg/backup
```

## Configuration

Initial configuration:
```
backup init
```
And fill in the rest.
