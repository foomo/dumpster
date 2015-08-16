package dumpster

import (
	"bytes"
	"errors"
	"os/exec"
)

// RestoreDump - restore a previous dump
func (d *Dumpster) RestoreDump(dumpType string, id string) (restoreReport string, restoreErrors string, err error) {
	dumpFilename, _ := d.getDumpFileNames(dumpType, id)
	dumpConfig, ok := d.Dumps[dumpType]
	if !ok {
		return "", "", errors.New("can not find this dump type: " + dumpType)
	}
	cmd := exec.Command(dumpConfig.Restore.Program, append(dumpConfig.Restore.Args, dumpFilename)...)
	bufferStdOut := bytes.NewBuffer([]byte{})
	bufferStdErr := bytes.NewBuffer([]byte{})
	cmd.Stderr = bufferStdErr
	cmd.Stdout = bufferStdOut
	exitError := cmd.Run()
	if exitError != nil {
		return "", "", errors.New("could not restore - calling restore program failed: " + exitError.Error())
	}
	return bufferStdOut.String(), bufferStdErr.String(), nil
}
