package tool

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// A Command defines the import path and provides a cached build.
type CommandBuild struct {
	ImportPath string
	buildPath  string
	buildOnce  sync.Once
	buildErr   error
}

// Build the Command and return the path to the binary.
func (c *CommandBuild) Build() (string, error) {
	c.buildOnce.Do(func() {
		basename := filepath.Base(c.ImportPath)
		exe, err := ioutil.TempFile("", basename+"-")
		if err != nil {
			c.buildErr = err
			return
		}
		c.buildPath = exe.Name()
		_ = os.Remove(c.buildPath) // the build tool will create this
		options := Options{
			ImportPaths: []string{c.ImportPath},
			Output:      c.buildPath,
		}
		if _, err = options.Command("build"); err != nil {
			c.buildErr = err
			return
		}
	})
	return c.buildPath, c.buildErr
}
