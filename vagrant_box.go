package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type VagrantBox struct {
	Name          string
	BaseDirectory string
	Directory     string
	Version       string
	Provider      string
	MasterID      string
}

func (b *VagrantBox) Update() error {
	if !b.IsExist() || !b.IsCloneMaster() {
		return nil
	}

	masterID, err := b.ReadMasterID()
	if err != nil {
		return err
	}
	b.MasterID = masterID

	return nil
}

func (b *VagrantBox) IsExist() bool {
	stat, err := os.Stat(b.Directory)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func (b *VagrantBox) IsCloneMaster() bool {
	stat, err := os.Stat(b.providerFile("master_id"))
	if err != nil {
		return false
	}
	return stat.Size() > 0
}

func (b *VagrantBox) ReadMasterID() (string, error) {
	buf, err := b.readProviderFile("master_id")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buf)), nil
}

func (b *VagrantBox) providerFile(filename string) string {
	return filepath.Join(b.Directory, filename)
}

func (b *VagrantBox) readProviderFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(b.providerFile(filename))
}
