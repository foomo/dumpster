package mock

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

// __DIR__
func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

// GetMockBytes mock data byte, tesing only - relative name from mock
func GetBytes(name string) []byte {
	f, err := os.Open(GetFilename(name))
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return b
}

func GetFilename(name string) string {
	return getCurrentDir() + "/" + name
}
