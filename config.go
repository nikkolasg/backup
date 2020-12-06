package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// BackupConfig contains all info
type BackupConfig struct {
	Remote    string
	LocalRoot string
	Upload    Config
	Download  Config
	// makes the endpoint the exact copy as the source
	Sync Config
}

const (
	Upload int = iota
	Download
	Sync
)

// WriteTo writes the config to the given writer
func (b *BackupConfig) WriteTo(w io.Writer) error {
	return toml.NewEncoder(w).Encode(b)
}

func (b BackupConfig) WriteToFile(fname string) error {
	return getWriter(fname, func(w io.Writer) error {
		return b.WriteTo(w)
	})
}

// Config specific for an action
type Config struct {
	Includes []string
	Excludes []string
}

// Add the given path to the corresponding type. Checks if file exists, when
// prefixed with the local root
func (b *BackupConfig) Add(t int, path string) error {
	if !fileExists(filepath.Join(b.LocalRoot, path)) {
		return fmt.Errorf("config: inexistant %s", path)
	}
	switch t {
	case Sync:
		b.Sync.Add(path)
	case Upload:
		b.Upload.Add(path)
	case Download:
		if !b.Upload.Contains(path) && !b.Sync.Contains(path) {
			return fmt.Errorf("adding path not existent in upload or sync list %s", path)
		}
		b.Download.Add(path)
	}
	return nil
}

// Add path to the config - path can be a file or a folder
func (c *Config) Add(path string) {
	c.Includes = append(c.Includes, path)
}

func (c *Config) Contains(path string) bool {
	for _, ipath := range c.Includes {
		if ipath == path {
			return true
		}
	}
	return false
}

// Load  returns a BackupConfig from a file
func Load(path string) (*BackupConfig, error) {
	var bc BackupConfig
	if _, err := toml.DecodeFile(path, &bc); err != nil {
		return nil, fmt.Errorf("error decoding config file: %s", err)
	}
	return &bc, nil
}

var dump = `
# Server where to backup/restore - one can write any information understandable 
# by rsync
#Remote = "foo@bar.com:/data/backup/"
# Root on local filesystem from where to base the rest of the paths
#LocalRoot = "/home/zoo/"

# Configuration items for the uploading part, i.e. the files or folder
# that you want to save on the remote server. Uploading de not delete
# any file on the remote server, it only overwrite existing file if 
# changed and copy the new files if any.
#[Upload]
  # List of folder/files to backup
  # Includes = ["prog/", "documents/", "movies/", "music/"]

  # List of folder/files to exclude from backup
  # It can be useful if there are some subfolder you wish to exclude
  # Excludes = ["prog/go/"]

# Configuration items for the downloading part, i.e. the files of folder
# that you want to restore from the remote server
# It contains the same fields than the Upload config (Include, IncludeFile...)
# Downloading do not delete any file locally, it only updates existing files
# if changed and copy new files if any.
# [Download]
  # Includes = ["documents/", "music/"]


# Configuration items for the files and folder you wish to sync
# Sync means it will recreate the exact content on the given paths
# locally or remotely depending on the action (upload or download)
# It can be useful for configuration folders for example that have
# new files overtime, and you don't want to keep the old files around,
# just want the latest version. Think of it as a git clone operation that 
# only takes the latest state.
# [Sync]
#   Includes = [".config/awesome/", ".config/termite/"]
`

var ConfigPerm os.FileMode = 0700

// DumpSampleConfig writes a sample config to the given path
func DumpSampleConfig(path string) error {
	return ioutil.WriteFile(path, []byte(dump), ConfigPerm)
}
