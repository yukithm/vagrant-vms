package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
)

type VirtualBoxVM struct {
	ID   string
	Name string
}

func GetVirtualBoxVMs() ([]VirtualBoxVM, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("VBoxManage", "list", "vms")
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var vms = make([]VirtualBoxVM, 0)
	re := regexp.MustCompile(`\A"([^"]+)"\s+\{([0-9a-f-]+)\}\z`)
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		m := re.FindStringSubmatch(scanner.Text())
		if m != nil {
			vm := VirtualBoxVM{
				ID:   m[2],
				Name: m[1],
			}
			vms = append(vms, vm)
		}
	}

	return vms, nil
}
