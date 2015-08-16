package dumpster

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

// RestoreRemoteDump restore a dump from a remote server
func (d *Dumpster) RestoreRemoteDump(remoteName string, dumpType string, id string) (restoreReport string, restoreErrors string, err error) {
	remote, ok := d.Remotes[remoteName]
	if !ok {
		return "", "", errors.New("unknown remote")
	}
	response, err := http.Get(remote.EndPoint + "/dump/" + url.QueryEscape(dumpType) + "/" + url.QueryEscape(id))
	if err != nil {
		return "", "", errors.New("could not get remote dump: " + err.Error())
	}
	defer response.Body.Close()
	tempFile, err := ioutil.TempFile(d.DataDir, "dumpster-restore-")
	if err != nil {
		return "", "", errors.New("could not create file to download remote dump: " + err.Error())
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		return "", "", err
	}
	restoreReport, restoreErrors, err = d.restoreDump(dumpType, tempFile.Name())
	err = os.Remove(tempFile.Name())
	if err != nil {
		err = errors.New("dump was restored, but the temp file could not be removed: " + err.Error())
	}
	log.Println("YEAH")
	return restoreReport, restoreErrors, err
}

// RestoreDump - restore a previous dump
func (d *Dumpster) RestoreDump(dumpType string, id string) (restoreReport string, restoreErrors string, err error) {
	dump, err := d.GetDump(dumpType, id)
	if err != nil {
		return "", "", err
	}
	return d.restoreDump(dumpType, dump.Filename)
}

func (d *Dumpster) restoreDump(dumpType string, dumpFilename string) (restoreReport string, restoreErrors string, err error) {
	dumpConfig, ok := d.Dumps[dumpType]
	if !ok {
		return "", "", errors.New("can not find this dump type: " + dumpType)
	}
	if dumpConfig.Restore.Program == "" {
		return "", "", errors.New("this dump can not be restored, because a restore program has not been configured")
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
