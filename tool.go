// Package tool wraps the go tool.
package tool

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Defines an Build command.
type Options struct {
	GoBin       string
	ImportPaths []string
	Output      string
	ForceAll    bool
	Parallel    uint
	Compiler    string
	GccGoFlags  string
	GcFlags     string
	LdFlags     string
	Tags        string
	Verbose     bool
}

type CommandError struct {
	command string
	options *Options
	out     []byte
	err     []byte
}

// Default fallback.
var goBinFallback string

func goBin(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if goBinFallback != "" {
		return goBinFallback, nil
	}
	var err error
	goBinFallback, err = exec.LookPath("go")
	if err != nil {
		return "", fmt.Errorf("Error finding go binary: %s", err)
	}
	return goBinFallback, nil
}

func (e *CommandError) Error() string {
	return fmt.Sprintf(
		"%s Failed %+v: %s\n%s", strings.Title(e.command), e.options, e.out, e.err)
}

func (e *CommandError) StdErr() []byte {
	return e.err
}

func (e *CommandError) StdOut() []byte {
	return e.out
}

func (o *Options) Command(command string) (affected []string, err error) {
	args := []string{command}
	if o.Output != "" {
		args = append(args, "-o", o.Output)
	}
	if o.ForceAll {
		args = append(args, "-a")
	}
	if o.Parallel != 0 {
		args = append(args, "-p", fmt.Sprintf("%d", o.Parallel))
	}
	if o.Compiler != "" {
		args = append(args, "-compiler", o.Compiler)
	}
	if o.GccGoFlags != "" {
		args = append(args, "-gccgoflags", o.GccGoFlags)
	}
	if o.GcFlags != "" {
		args = append(args, "-gcflags", o.GcFlags)
	}
	if o.LdFlags != "" {
		args = append(args, "-ldflags", o.LdFlags)
	}
	if o.Tags != "" {
		args = append(args, "-tags", o.Tags)
	}
	if o.Verbose {
		args = append(args, "-v")
	}
	for _, importPath := range o.ImportPaths {
		args = append(args, importPath)
	}
	bin, err := goBin(o.GoBin)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(bin, args...)
	var bufOut, bufErr bytes.Buffer
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr
	err = cmd.Run()
	if err != nil {
		return nil, &CommandError{
			command: command,
			options: o,
			out:     bufOut.Bytes(),
			err:     bufErr.Bytes(),
		}
	}
	affectedBytes := bytes.Split(bufErr.Bytes(), []byte("\n"))
	affected = make([]string, 0, len(affectedBytes))
	for _, importPath := range affectedBytes {
		if len(importPath) == 0 {
			continue
		}
		affected = append(affected, string(importPath))
	}
	return affected, nil
}
