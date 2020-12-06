package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

var defaultFolder = filepath.Join(os.Getenv("HOME"), ".config/backup")
var defaultConfig = "backup.toml"

var confFlag = &cli.StringFlag{
	Name:  "c, config",
	Usage: "Load configuration from `FOLDER` - default is $HOME/.config/backup/",
}

var fakeFlag = &cli.BoolFlag{
	Name:  "f, fake",
	Usage: "only prints command",
}

var verboseFlag = &cli.BoolFlag{
	Name:  "v, verbose",
	Usage: "verbose output",
}

var noSyncFlag = &cli.BoolFlag{
	Name:  "nosync",
	Usage: "avoid sync operations (useful for quick upload)",
}

var syncFlag = &cli.BoolFlag{
	Name:  "sync",
	Usage: "Mark the upload folder to sync - it only makes sense for folders",
}

var fromFlag = &cli.StringFlag{
	Name:  "from",
	Usage: "Fetch configuration from `USER@HOST`",
}

var cleanFlag = &cli.BoolFlag{
	Name:  "clean",
	Usage: "clean the existing config",
}

var remoteFlag = &cli.StringFlag{
	Name:     "remote",
	Usage:    "remote endpoint - rsync must be able to parse it",
	Required: true,
}

var localFlag = &cli.StringFlag{
	Name:     "local",
	Usage:    "local base folder - rsync ",
	Required: true,
}

var app = &cli.App{
	Commands: []cli.Command{
		{
			Name:  "upload",
			Usage: "upload to remote host",
			Flags: []cli.Flag{confFlag, fakeFlag, verboseFlag, noSyncFlag},
			Action: func(c *cli.Context) error {
				banner(c)
				return upload(c)
			},
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "adds `PATH` `PATH2` ... to the upload list",
					Flags: []cli.Flag{confFlag, verboseFlag, syncFlag},
					Action: func(c *cli.Context) error {
						banner(c)
						return uploadAdd(c)
					},
				},
			},
		},
		{
			Name:  "download",
			Usage: "download from a remote host",
			Flags: []cli.Flag{confFlag, fakeFlag, verboseFlag},
			Action: func(c *cli.Context) error {
				banner(c)
				return downloadCmd(c)
			},
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "adds `PATH` to the download list",
					Flags: []cli.Flag{confFlag, verboseFlag},
					Action: func(c *cli.Context) error {
						banner(c)
						return downloadAdd(c)
					},
				},
			},
		},
		{
			Name:  "init",
			Usage: "Init a backup configuration or fetch from server",
			Flags: []cli.Flag{confFlag, cleanFlag, localFlag, remoteFlag},
			Action: func(c *cli.Context) error {
				banner(c)
				return initConfig(c)
			},
		},
	},
}

func main() {
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
	fmt.Printf("Uploading to %s\n", bc.Remote)
	rsync := newRsync(bc.LocalRoot, bc.Remote)
	defer rsync.Cleanup()
	err = rsync.Upload(bc.Upload)
	return err
}

func uploadAdd(c *cli.Context) error {
	bc, err := getConf(c)
	if err != nil {
		return err
	}
	for i := 0; i < c.NArg(); i++ {
		var err error
		path := c.Args().Get(i)
		err = bc.Add(Upload, path)
		if err != nil {
			return fmt.Errorf("upload add: err adding %s: %s", path, err)
		}
	}
	folder := getConfFolder(c)
	fname := getConfigFile(folder)
	return bc.WriteToFile(fname)
}

func downloadAdd(c *cli.Context) error {
	bc, err := getConf(c)
	if err != nil {
		return err
	}
	path := c.Args().Get(0)
	if err := bc.Add(Download, path); err != nil {
		return fmt.Errorf("error adding download: %s", err)
	}
	folder := getConfFolder(c)
	fname := getConfigFile(folder)
	return bc.WriteToFile(fname)
}

func downloadCmd(c *cli.Context) error {
	bc, err := getConf(c)
	if err != nil {
		return err
	}
	rsync := newRsync(bc.LocalRoot, bc.Remote)
	defer rsync.Cleanup()
	return rsync.Download(bc.Download)
}

func initConfig(c *cli.Context) error {
	if c.IsSet(fromFlag.Name) {
		//return fetchConfig(c.String(fromFlag.Name))
	}

	folder := getConfFolder(c)
	if err := os.MkdirAll(folder, 0700); err != nil {
		return err
	}

	fname := getConfigFile(folder)
	if c.Bool(cleanFlag.Name) {
		fmt.Printf("Confirmation of resetting config file (type YES):\n")
		if !askConfirmation() {
			return nil
		}
		deleteFile(fname)
	}
	var backup BackupConfig
	backup.Remote = c.String(remoteFlag.Name)
	backup.LocalRoot = c.String(localFlag.Name)
	if err := backup.WriteToFile(fname); err != nil {
		return fmt.Errorf("error writing config: %s", err)
	}
	fmt.Println("Config file written at ", fname)
	return nil
}

func getConf(c *cli.Context) (*BackupConfig, error) {
	folder := getConfFolder(c)
	path := getConfigFile(folder)
	return Load(path)
}

func getConfigFile(folder string) string {
	return filepath.Join(folder, defaultConfig)
}

func getConfFolder(c *cli.Context) string {
	if c.IsSet("config") {
		return c.String("config")
	}
	return defaultFolder
}

func askConfirmation() bool {
	in := readOneLine()
	return strings.Contains(in, "YES")
}

var input io.Reader = os.Stdin

func readOneLine() string {
	reader := bufio.NewReader(input)
	s, _ := reader.ReadString('\n')
	return s
}
