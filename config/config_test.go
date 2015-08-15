package config

import (
	"testing"

	"github.com/foomo/dumpster/mock"
	"github.com/foomo/dumpster/utils"
)

func TestConfig(t *testing.T) {

	c, err := LoadConfig(mock.GetFilename("config.yml"))
	utils.PanicOnErr(err)
	if c.Remotes["hansi"].EndPoint != "http://hansi/" {
		t.Fatal("could not read remote hansi:", c.Remotes["hansi"].EndPoint)
	}
}
