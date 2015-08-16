package dumpster

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/foomo/dumpster/config"
)

const metaPrefix = ".meta-"

type Dump struct {
	ID           string
	Created      time.Time
	Comment      string
	Report       string
	Errors       string
	Type         string
	Filename     string
	MetaFilename string
}

type Dumpster struct {
	DataDir string
	Dumps   map[string]config.Dump
	Remotes map[string]config.Remote
}

var (
	ErrNotFound = errors.New("dump not found")
)

func NewDumpster(dataDir string, dumps map[string]config.Dump, remotes map[string]config.Remote) (d *Dumpster, err error) {
	d = &Dumpster{
		DataDir: dataDir,
		Dumps:   dumps,
		Remotes: remotes,
	}
	// try to create dump dirs
	for dumpType := range d.Dumps {
		dumpTypeValidationErr := validatePathPart(dumpType)
		if dumpTypeValidationErr != nil {
			return nil, errors.New("invalid dump type: " + dumpType + " error: " + dumpTypeValidationErr.Error())
		}
		dumpDir := d.getDumpDir(dumpType)
		err := os.MkdirAll(dumpDir, 0777)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (d *Dumpster) GetDumps() map[string][]*Dump {
	dumps := make(map[string][]*Dump)
	for dumpType := range d.Dumps {
		dumps[dumpType] = d.GetDumpsForType(dumpType)
	}
	return dumps
}

func (d *Dumpster) GetDumpsForType(dumpType string) []*Dump {
	dumpDir := d.getDumpDir(dumpType)
	files, err := ioutil.ReadDir(dumpDir)
	dumps := []*Dump{}
	if err == nil {
		for _, file := range files {
			if strings.HasPrefix(file.Name(), metaPrefix) {
				metaFilename := path.Join(dumpDir, file.Name())
				metaBytes, err := ioutil.ReadFile(metaFilename)
				if err == nil {
					dump := &Dump{}
					jsonErr := json.Unmarshal(metaBytes, &dump)
					if jsonErr == nil {
						dumpFilename, _ := d.getDumpFileNames(dumpType, dump.ID)
						dump.Filename = dumpFilename
						dump.MetaFilename = metaFilename
						dump.Type = dumpType
						dumps = append(dumps, dump)
					}
				}
			}
		}
	}
	return dumps
}

func (d *Dumpster) GetDump(dumpType, id string) (dump *Dump, err error) {
	for _, dump := range d.GetDumpsForType(dumpType) {
		if dump.ID == id {
			return dump, nil
		}
	}
	return nil, ErrNotFound
}

func (d *Dumpster) DeleteDump(dumpType, id string) error {
	dump, err := d.GetDump(dumpType, id)
	if err != nil {
		return err
	}
	err = os.Remove(dump.MetaFilename)
	if err != nil {
		return errors.New("could not remove meta file")
	}
	err = os.Remove(dump.Filename)
	if err != nil {
		return errors.New("could not remove dump file")
	}
	return nil
}

func (d *Dumpster) getDumpFileNames(dumpType string, id string) (dumpFile string, dumpMetaDataFile string) {
	base := d.getDumpDir(dumpType)
	return path.Join(base, id), path.Join(base, metaPrefix+id)
}

func (d *Dumpster) getDumpDir(dumpType string) string {
	return path.Join(d.DataDir, dumpType)
}

const pathRuleRegex = `^([a-z]|[A-Z]|[0-9]|[\-_])+$`

func validatePathPart(part string) error {
	var validateID = regexp.MustCompile(pathRuleRegex)
	if !validateID.MatchString(part) {
		return errors.New("forbidden input - does not match: " + pathRuleRegex)
	}
	return nil
}
