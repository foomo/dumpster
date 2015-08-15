package dumpster

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"time"
)

// CreateDump creates a dump, by calling the configured dump program
func (d *Dumpster) CreateDump(dumpType string, id string, comment string) (metaData *Dump, err error) {
	dumpFilename, metaDataFilename := d.getDumpFileNames(dumpType, id)
	dumpConfig, ok := d.Dumps[dumpType]
	if !ok {
		return nil, errors.New("can not find this dump type: " + dumpType)
	}
	cmd := exec.Command(dumpConfig.Dump.Program, append(dumpConfig.Dump.Args, dumpFilename)...)
	bufferStdOut := bytes.NewBuffer([]byte{})
	bufferStdErr := bytes.NewBuffer([]byte{})
	cmd.Stderr = bufferStdErr
	cmd.Stdout = bufferStdOut
	exitError := cmd.Run()
	if exitError != nil {
		return nil, errors.New("could not dump - calling dump program failed: " + exitError.Error())
	}
	metaData = &Dump{
		ID:      id,
		Comment: comment,
		Created: time.Now(),
		Report:  string(bufferStdOut.Bytes()),
		Errors:  string(bufferStdErr.Bytes()),
		Type:    dumpType,
	}
	metaDataBytes, err := json.MarshalIndent(metaData, "", "    ")
	if err != nil {
		return nil, errors.New("could not marshal metadata: " + err.Error())
	}
	err = ioutil.WriteFile(metaDataFilename, metaDataBytes, 0644)
	if err != nil {
		return nil, err
	}
	return metaData, nil
}
