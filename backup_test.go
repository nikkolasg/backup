package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	tmp := tmpFolder()
	defer os.RemoveAll(tmp)
	defaultFolder = tmp

	args := toArgs("init", "--config", tmp, "--remote", tmp, "--local", tmp)
	require.NoError(t, app.Run(args))
	defer func() {
		input = os.Stdin
	}()
	var in bytes.Buffer
	in.WriteString("YES\n")
	input = &in
	toSearch := "BROU"
	args = toArgs("init", "--config", tmp, "--clean", "--remote", toSearch, "--local", tmp)
	require.NoError(t, app.Run(args))

	bc, err := Load(getConfigFile(tmp))
	require.NoError(t, err)
	require.Equal(t, bc.LocalRoot, tmp)
	require.Equal(t, bc.Remote, toSearch)
}

func TestBackup(t *testing.T) {
	source := tmpFolder()
	defer os.RemoveAll(source)
	dest := tmpFolder()
	defer os.RemoveAll(dest)
	args := toArgs("init", "--config", source, "--remote", dest, "--local", source)
	require.NoError(t, app.Run(args))

	fname1 := "test.md"
	createTestFile(filepath.Join(source, fname1))
	args = toArgs("upload", "add", "--config", source, fname1)
	require.NoError(t, app.Run(args))

	loadAndCheck := func(fn func(*BackupConfig) bool) {
		bc, err := Load(getConfigFile(source))
		require.NoError(t, err)
		require.NotNil(t, bc)
		require.True(t, fn(bc))
	}

	loadAndCheck(func(bc *BackupConfig) bool { return bc.Upload.Contains(fname1) })

	fname2 := "test2.md"
	createTestFile(filepath.Join(source, fname2))
	args = toArgs("upload", "add", "--config", source, "--sync", fname2)
	require.NoError(t, app.Run(args))
	loadAndCheck(func(bc *BackupConfig) bool { return bc.Sync.Contains(fname2) })

	args = toArgs("download", "add", "--config", source, fname2)
	require.NoError(t, app.Run(args))
	loadAndCheck(func(bc *BackupConfig) bool { return bc.Download.Contains(fname2) })

	args = toArgs("download", "add", "--config", source, "bloup")
	require.Error(t, app.Run(args))

	// test to upload!
	args = toArgs("upload", "--config", source)
	require.NoError(t, app.Run(args))
	full1 := filepath.Join(dest, fname1)
	full2 := filepath.Join(dest, fname2)
	require.True(t, fileExists(full1))
	require.True(t, fileExists(full2))
	// test to download
	deleteFile(full2)
	args = toArgs("download", "--config", source)
	require.NoError(t, app.Run(args))
	require.True(t, fileExists(full2))

}

func TestGenerateDump(t *testing.T) {
	var bc BackupConfig
	bc.Remote = "foo@bar:/media/local"
	bc.LocalRoot = "/home/local"
	bc.Upload.Includes = []string{"/my/file1", "/my/top/folder"}
	bc.Upload.Excludes = []string{"/my/top/folder/cat.jpg"}
	//bc.Download.Includes = []string{"/my}
}

func toArgs(args ...string) []string {
	return append([]string{"backup"}, args...)
}

func createTestFile(path string) {
	ioutil.WriteFile(path, []byte(""), 0700)
}
