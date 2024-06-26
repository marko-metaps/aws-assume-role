package main

import (
	"fmt"
	"os"
)

type FileSystem interface {
	UserHomeDir() (string, error)
	Stat(path string) (os.FileInfo, error)
}

type RealFileSystem struct{}

func NewRealFileSystem() *RealFileSystem {
	return &RealFileSystem{}
}

func (RealFileSystem) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (RealFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func checkCredentialFile(fs FileSystem) error {
	dir, err := fs.UserHomeDir()
	if err != nil {
		return err
	}
	path := dir + "/.aws/credentials"
	stat, err := fs.Stat(path)
	if os.IsNotExist(err) {
		return err
	}
	if stat.Mode().Perm()&0200 != 0200 {
		return fmt.Errorf("no write permission. [~/.aws/credentials]")
	}
	return nil
}
