package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// BackupConfig contains all info
type BackupConfig struct {
	Remote     string
	RemoteRoot string
	LocalRoot  string
	Upload     Config
	Download   Config
	// makes the endpoint the exact copy as the source
	Sync Config
}

// Config specific for an action
type Config struct {
	Includes []string
	Excludes []string
}

// Load  returns a BackupConfig from a file
func Load(path string) (*BackupConfig, error) {
	var bc BackupConfig
	if _, err := toml.DecodeFile(path, &bc); err != nil {
		return nil, fmt.Errorf("error decoding config file: %s", err)
	}
	return &bc, nil
}

// DumpSampleConfig writes a sample config to the given path
func DumpSampleConfig(path string) {
	conf := &BackupConfig{
		Remote:    "foo@bar.com:/data/backup/",
		LocalRoot: "/home/zoo/",
		Upload: Config{
			Includes: []string{"prog/", "documents/", "movies/", "music/"},
			Excludes: []string{"prog/go/"},
		},
		Download: Config{
			Includes: []string{"documents/", "music/"},
		},
		Sync: Config{
			Includes: []string{".config/awesome/", ".config/termite/"},
		},
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("error creating dump config: %s", err)
		return
	}
	if err := toml.NewEncoder(f).Encode(&conf); err != nil {
		fmt.Printf("error encoding file: %s", err)
	}
}
