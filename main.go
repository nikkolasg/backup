package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var confFlag = &cli.StringFlag{
	Name:     "c, config",
	Usage:    "Load configuration from `FILE`, REQUIRED.",
	Required: true,
}

var fakeFlag = &cli.BoolFlag{
	Name:  "f, fake",
	Usage: "only prints command",
}

var verboseFlag = &cli.BoolFlag{
	Name:  "v, verbose",
	Usage: "verbose output",
}

func main() {
	app := &cli.App{
		Commands: []cli.Command{
			{
				Name:  "upload",
				Usage: "upload to remote host",
				Flags: []cli.Flag{confFlag, fakeFlag, verboseFlag},
				Action: func(c *cli.Context) error {
					banner(c)
					return upload(c)
				},
			},
			{
				Name:  "download",
				Usage: "download from a remote host",
				Flags: []cli.Flag{confFlag, fakeFlag, verboseFlag},
				Action: func(c *cli.Context) error {
					banner(c)
					return nil
				},
			},
			{
				Name:  "init",
				Usage: "fetch config file from server",
				Action: func(c *cli.Context) error {
					banner(c)
					panic("not implemented")
				},
			},
			{
				Name:  "example",
				Usage: "writes an example config",
				Action: func(c *cli.Context) error {
					banner(c)
					return example(c)
				},
			},
		},
	}
	app.Run(os.Args)
}

func banner(c *cli.Context) {
	Fake = c.Bool("fake")
	Verbose = c.Bool("verbose")
	if Verbose {
		fmt.Printf("backup tool. options: fake(%v), verbose(%v)\n", Fake, Verbose)
	}

}

func handle(err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
	return nil
}

func upload(c *cli.Context) error {
	bc, err := getConf(c)
	if err != nil {
		return err
	}
	fmt.Printf("Uploading to %s:%s\n", bc.Remote, bc.RemoteRoot)
	rsync := newRsync(bc.LocalRoot, bc.Remote, bc.RemoteRoot)
	return rsync.Upload(bc.Upload.Includes, bc.Upload.Excludes)
}

func example(c *cli.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fname := filepath.Join(dir, "example.toml")
	DumpSampleConfig(fname)
	fmt.Printf("- Wrote example configuration in %s\n", fname)
	return nil
}

func getConf(c *cli.Context) (*BackupConfig, error) {
	path := c.String("config")
	conf, err := Load(path)
	return conf, err
}
