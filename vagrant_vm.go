package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// VagrantVM represents vagrant's VM information.
type VagrantVM struct {
	ID        string
	ShortID   string
	Name      string
	Provider  string
	State     string
	Directory string
	VMID      string
}

func (vm *VagrantVM) Update() error {
	if !vm.IsExist() {
		return nil
	}

	id, err := vm.ReadIndexUUID()
	if err != nil {
		return err
	}
	vm.ID = id

	vmid, err := vm.ReadVMID()
	if err != nil {
		return err
	}
	vm.VMID = vmid

	return nil
}

func (vm *VagrantVM) IsExist() bool {
	stat, err := os.Stat(vm.ProviderDirectory())
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func (vm *VagrantVM) ProviderDirectory() string {
	return filepath.Join(vm.Directory, ".vagrant", "machines", vm.Name, vm.Provider)
}

func (vm *VagrantVM) ReadIndexUUID() (string, error) {
	buf, err := vm.readProviderFile("index_uuid")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buf)), nil
}

func (vm *VagrantVM) ReadVMID() (string, error) {
	buf, err := vm.readProviderFile("id")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buf)), nil
}

func (vm *VagrantVM) providerFile(filename string) string {
	return filepath.Join(vm.ProviderDirectory(), filename)
}

func (vm *VagrantVM) readProviderFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(vm.providerFile(filename))
}

func GlobalStatus() ([]VagrantVM, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("vagrant", "global-status")
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var vms = make([]VagrantVM, 0)
	re := regexp.MustCompile(`\A([0-9a-f]+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(.+?)\s*\z`)
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		m := re.FindStringSubmatch(scanner.Text())
		if m != nil {
			vm := VagrantVM{
				ID:        m[1],
				ShortID:   m[1],
				Name:      m[2],
				Provider:  m[3],
				State:     m[4],
				Directory: m[5],
			}
			if vm.IsExist() {
				if err := vm.Update(); err != nil {
					return nil, err
				}
			}
			vms = append(vms, vm)
		}
	}
	return vms, nil
}

func findVagrantCommand() (string, error) {
	path, err := exec.LookPath("vagrant")
	if err != nil {
		return "", err
	}
	return path, nil
}
