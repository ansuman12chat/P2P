package app

import (
	"os"
)

type Ioutiler interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type Ioutil struct{}

func (a Ioutil) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (a Ioutil) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
