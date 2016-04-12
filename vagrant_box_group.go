package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type VagrantBoxGroup struct {
	Name      string
	Directory string
	Boxes     []VagrantBox
}

func (bg *VagrantBoxGroup) Add(box VagrantBox) {
	for _, b := range bg.Boxes {
		if boxEqual(b, box) {
			return
		}
	}
	bg.Boxes = append(bg.Boxes, box)
}

func (bg *VagrantBoxGroup) Versions() []string {
	versions := make([]string, 0, len(bg.Boxes))
	for _, box := range bg.Boxes {
		versions = append(versions, box.Version)
	}
	return versions
}

func GetVagrantBoxes(directory string) ([]VagrantBoxGroup, error) {
	var boxes = make([]VagrantBoxGroup, 0)
	err := eachFiles(directory, func(file os.FileInfo, path string) error {
		box, err := getVagrantBoxGroup(path)
		if err != nil {
			return err
		}
		boxes = append(boxes, box)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

func getVagrantBoxGroup(boxDir string) (VagrantBoxGroup, error) {
	boxName := toBoxName(filepath.Base(boxDir))
	group := VagrantBoxGroup{
		Name:      boxName,
		Directory: boxDir,
	}

	err := eachFiles(boxDir, func(file os.FileInfo, path string) error {
		if !file.IsDir() {
			return nil
		}
		boxes, err := getProviderBoxes(path, boxName)
		if err != nil {
			return err
		}
		for _, box := range boxes {
			group.Add(box)
		}
		return nil
	})
	if err != nil {
		return VagrantBoxGroup{}, err
	}

	return group, nil

}

func getProviderBoxes(versionDir string, boxName string) ([]VagrantBox, error) {
	var boxes = make([]VagrantBox, 0)
	err := eachFiles(versionDir, func(file os.FileInfo, path string) error {
		if !file.IsDir() {
			return nil
		}

		box := VagrantBox{
			Name:          boxName,
			BaseDirectory: filepath.Dir(versionDir),
			Directory:     path,
			Version:       filepath.Base(versionDir),
			Provider:      file.Name(),
		}
		if err := box.Update(); err != nil {
			return err
		}

		boxes = append(boxes, box)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

func eachFiles(dir string, f func(os.FileInfo, string) error) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if err := f(file, path); err != nil {
			return err
		}
	}

	return nil
}

func toBoxName(dirname string) string {
	return strings.Replace(dirname, "-VAGRANTSLASH-", "/", -1)
}

func boxEqual(a, b VagrantBox) bool {
	return a.Name == b.Name && a.Version == b.Version && a.Provider == b.Provider
}
