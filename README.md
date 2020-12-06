[![Build Status](https://travis-ci.com/nikkolasg/backup.svg?branch=master)](https://travis-ci.com/nikkolasg/backup)
[![codecov](https://codecov.io/gh/nikkolasg/backup/branch/master/graph/badge.svg?token=L0JCS6XYXH)](https://codecov.io/gh/nikkolasg/backup)

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

backup requires to know at least two informations:
* Local root: the absolute path from which where to backup files and folders. It
  can be `/home/ubuntu` for example, or `/` if you wish to backup from your root
  directory. This is required from rsync.
* Remote root: the absolute path on the destination. This must be understood as
  rsync so it can be `/backup/ubuntu` or `ubuntu@192.192.192.192:/backup/ubuntu`
  for a remote server destination.

```
backup init --remote <remote> --local <local>
```

If you want to use a custom configuration folder

#### Adding folders to upload

Adding folders / files is as easy as:
```
backup upload add <path1> <path2> ...
```

TODO: add `--download` to automatically add it to download.

#### Adding folders to download

You must specify all folders to download to backup.
You don't necessarily want to download all files you have backed up (for example
large files that you are fine keeping on the destination path).
```
backup download add <path1> ...
```

Note that `path1` **must** be included in the upload list to be valid. `backup`
doesn't allow to restore a file not tracked in upload yet.

TODO: make list argument

TODO: git restore
