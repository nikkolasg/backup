package main

import (
	"fmt"
	"io"
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
}

const (
	Upload int = iota
	Download
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
	case Upload:
		b.Upload.Add(path)
	case Download:
		if !b.Upload.Contains(path) {
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

var ConfigPerm os.FileMode = 0700
