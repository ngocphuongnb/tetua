package test

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetRootDir return the root directory of the project
func GetRootDir() string {
	_, mainFile, _, _ := runtime.Caller(2)
	return filepath.Dir(mainFile)
}

// CreateDir creates a temporary directory for testing
func CreateDir(prefix string) string {
	tmpDir := path.Join(path.Dir(path.Dir(GetRootDir())), "private/tmp")
	dir, err := ioutil.TempDir(tmpDir, prefix)
	if err != nil {
		panic(err)
	}
	return dir
}

func RecoverPanic(t *testing.T, expected string, msgs ...string) {
	msg := ""
	if len(msgs) > 0 {
		msg = "should panic on " + msgs[0]
	}
	if rs := recover(); rs != nil {
		if err, ok := rs.(error); ok && err != nil {
			assert.Equal(t, expected, err.Error(), msg)
		} else {
			assert.Equal(t, expected, rs, msg)
		}
	}
}
