package config

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Dump struct {
	Dump    string
	Restore string
}

type Remote struct {
	EndPoint string "endpoint"
}

type Config struct {
	Address string
	DataDir string "datadir"
	Dumps   map[string]Dump
	Remotes map[string]Remote
}

func LoadConfig(file string) (c *Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	yamlBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	c = &Config{}
	err = yaml.Unmarshal(yamlBytes, c)
	if err != nil {
		return nil, errors.New("could not parse config:" + err.Error())
	}
	return c, nil
}
