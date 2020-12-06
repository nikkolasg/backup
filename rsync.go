package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Fake only prints the commands out
var Fake = false

// Verbose print outputs of the commands
var Verbose = false

type rsync struct {
	rhost string // user@host:/data/backup
	lroot string // /home/user/ ..
	tmp   string
}

func newRsync(lroot, rhost string) *rsync {

	return &rsync{
		rhost: rhost,
		lroot: lroot,
		tmp:   tmpFolder(),
	}
}

func (r *rsync) Download(download Config) error {
	include, exclude, err := r.getFilePath(download)
	if err != nil {
		return err
	}
	cmd := r.baseCmd(include, exclude)
	cmd = r.toUpload(cmd)
	return run(cmd)

}

func (r *rsync) Cleanup() {
	os.RemoveAll(r.tmp)
}

func (r *rsync) Upload(upload Config) error {
	include, exclude, err := r.getFilePath(upload)
	if err != nil {
		return err
	}
	cmd := r.baseCmd(include, exclude)
	cmd = r.toUpload(cmd)
	return run(cmd)
}

func (r *rsync) getFilePath(c Config) (include string, exclude string, err error) {
	if len(c.Includes) == 0 {
		return "", "", fmt.Errorf("Invalid config: no include or includefile")
	}

	var includePath, excludePath string
	includePath = filepath.Join(r.tmp, "include.backup")
	if err := writeTmp(includePath, c.Includes); err != nil {
		return "", "", err
	}
	if len(c.Excludes) > 0 {
		excludePath = filepath.Join(r.tmp, "exclude.backup")
		if err := writeTmp(excludePath, c.Excludes); err != nil {
			return "", "", err
		}
	}
	return includePath, excludePath, nil
}

func (r *rsync) SyncUpload(sync Config) error {
	cmd, err := r.syncCmd(sync)
	if err != nil {
		return err
	}
	cmd = r.toUpload(cmd)
	return run(cmd)
}

func (r *rsync) SyncDownload(sync Config) error {
	cmd, err := r.syncCmd(sync)
	if err != nil {
		return err
	}
	cmd = r.toDownload(cmd)
	return run(cmd)

}

func (r *rsync) toUpload(cmd []string) []string {
	cmd = append(cmd, r.lroot)
	cmd = append(cmd, r.rhost)
	return cmd
}

func (r *rsync) toDownload(cmd []string) []string {
	cmd = append(cmd, r.rhost)
	cmd = append(cmd, r.lroot)
	return cmd
}

func (r *rsync) syncCmd(sync Config) ([]string, error) {
	include, _, err := r.getFilePath(sync)
	if err != nil {
		return nil, err
	}
	cmd := r.baseCmd(include, "")
	cmd = append(cmd, "--delete")
	return cmd, nil
}

func writeTmp(fname string, paths []string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, p := range paths {
		f.WriteString(p + "\n")
	}
	return nil
}

func (r *rsync) baseCmd(include, exclude string) []string {
	cmd := []string{"rsync", "-ravzz", "--links", "--progress"}
	cmd = append(cmd, []string{"--files-from", include}...)
	if exclude != "" {
		cmd = append(cmd, []string{"--exclude-from", exclude}...)
	}
	return cmd
}

func preprocess(paths []string, fn func(p string) (string, string)) []string {
	prefixed := make([]string, 0, len(paths)*2)
	for _, p := range paths {
		flag, value := fn(p)
		prefixed = append(prefixed, flag)
		prefixed = append(prefixed, value)
	}
	return prefixed
}

func run(c []string) error {
	if Verbose {
		fmt.Println("command: " + strings.Join(c, " "))
	}
	if Fake {
		cmd := exec.Command(c[0], append([]string{"--dry-run"}, c[1:]...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return nil
	}
	cmd := exec.Command(c[0], c[1:]...)
	var err error
	var out string
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	} else {
		var buff []byte
		buff, err = cmd.CombinedOutput()
		out = string(buff)
	}
	if err != nil {
		return fmt.Errorf("error executing command:\n\t-%v\n\t-Error: %s\n\t-Output: %s", c, err, string(out))
	}
	return nil
}
