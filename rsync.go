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
	rhost string // user@host
	lroot string // /home/user/ ..
	rroot string // /data/backup/...
}

func newRsync(lroot, rhost, rroot string) *rsync {
	return &rsync{
		rhost: rhost,
		lroot: lroot,
		rroot: rroot,
	}
}

func (r *rsync) Download(paths, excludes []string) error {
	cmd := r.baseCmd(paths, excludes)
	// from remote host to local
	cmd = append(cmd, r.remoteEndpoint())
	cmd = append(cmd, r.lroot)
	return run(cmd)
}

func (r *rsync) Upload(paths, excludes []string) error {
	cmd := r.baseCmd(paths, excludes)
	// from local to remote
	cmd = append(cmd, r.lroot)
	cmd = append(cmd, r.remoteEndpoint())
	return run(cmd)
}

func (r *rsync) remoteEndpoint() string {
	return r.rhost + ":" + r.rroot
}

func (r *rsync) baseCmd(paths, excludes []string) []string {
	inclusion := preprocess(paths, func(p string) string {
		return "--include='" + filepath.Join(p, "***") + "'"
	})
	exclusion := preprocess(excludes, func(p string) string {
		return "--exclude='" + p + "'"
	})
	cmd := []string{"rsync", "-ravzz", "--links", "--progress"}
	cmd = append(cmd, inclusion...)
	cmd = append(cmd, exclusion...)
	cmd = append(cmd, "--exclude='*'")
	return cmd
}

func preprocess(paths []string, fn func(p string) string) []string {
	prefixed := make([]string, len(paths))
	copy(prefixed, paths)
	for i, p := range prefixed {
		prefixed[i] = fn(p)
	}
	return prefixed
}

func run(c []string) error {
	if Verbose {
		fmt.Println("command: " + strings.Join(c, " "))
	}
	if Fake {
		/*cmd := exec.Command(c[0], append([]string{"--dry-run"}, c[1:]...)...)*/
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		/*cmd.Run()*/
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
