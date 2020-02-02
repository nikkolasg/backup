package main

import (
	"fmt"
	"io/ioutil"
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
	tmp, err := ioutil.TempDir("", "backup")
	if err != nil {
		panic(err)
	}
	return &rsync{
		rhost: rhost,
		lroot: lroot,
		tmp:   tmp,
	}
}

func (r *rsync) Download(paths, excludes []string) error {
	panic("not implemented yet")
	/*cmd := r.baseCmd(paths, excludes)*/
	//// from remote host to local
	//cmd = append(cmd, r.remoteEndpoint())
	//cmd = append(cmd, r.lroot)
	/*return run(cmd)*/
}

func (r *rsync) Cleanup() {
	os.RemoveAll(r.tmp)
}

func (r *rsync) Upload(paths, excludes []string) error {
	inclusion := filepath.Join(r.tmp, "include.backup")
	if err := writeTmp(inclusion, paths); err != nil {
		return err
	}

	exclusion := filepath.Join(r.tmp, "exclude.backup")
	if err := writeTmp(exclusion, excludes); err != nil {
		return err
	}

	cmd := r.baseCmd(inclusion, exclusion)
	// from local to remote
	cmd = append(cmd, r.lroot)
	cmd = append(cmd, r.rhost)
	return run(cmd)
}

func (r *rsync) SyncUpload(paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	cmd, err := r.syncCmd(paths)
	if err != nil {
		return err
	}

	cmd = append(cmd, r.lroot)
	cmd = append(cmd, r.rhost)
	return run(cmd)
}

func (r *rsync) syncCmd(paths []string) ([]string, error) {
	inclusion := filepath.Join(r.tmp, "include.backup")
	if err := writeTmp(inclusion, paths); err != nil {
		return nil, err
	}
	cmd := r.baseCmd(inclusion, "")
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
