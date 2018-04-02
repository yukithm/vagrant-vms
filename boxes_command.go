package main

import (
	"html/template"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var boxesDetailTemplate = `
{{- range $i, $e := . -}}
--[ {{add $i 1}} ]-------------------------------------
     name: {{.Name}}
directory: {{.Directory}}
 variants:{{range $bi, $e := .Boxes}}
      version: {{.Version}}
     provider: {{.Provider}}
    master_id: {{.MasterID}}
{{end}}
{{end -}}
`

func boxesCommand(c *cli.Context) {
	groups, err := GetVagrantBoxes(GetVagrantBoxDir())
	if err != nil {
		log.Fatal(err)
	}
	if c.Bool("vertical") {
		printBoxesVertically(groups)
	} else {
		printBoxes(groups)
	}
}

func maxVersionLen(groups []VagrantBoxGroup) (max int) {
	for _, group := range groups {
		for _, box := range group.Boxes {
			l := len(box.Version)
			if l > max {
				max = l
			}
		}
	}
	return
}

func printBoxes(groups []VagrantBoxGroup) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetHeaderLine(true)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"name", "version", "provider", "directory", "master_id"})
	for _, group := range groups {
		for _, box := range group.Boxes {
			table.Append([]string{
				box.Name,
				box.Version,
				box.Provider,
				box.BaseDirectory,
				box.MasterID,
			})
		}
	}
	table.Render()
}

func printBoxesVertically(groups []VagrantBoxGroup) {
	tmpl := template.New("boxes").Funcs(template.FuncMap{
		"add": tmplFuncAdd,
	})
	tmpl = template.Must(tmpl.Parse(boxesDetailTemplate))
	err := tmpl.Execute(os.Stdout, groups)
	if err != nil {
		log.Fatal(err)
	}
}

// GetVagrantBoxDir returns vagrant's boxes directory.
func GetVagrantBoxDir() string {
	dir, err := homedir.Expand("~/.vagrant.d/boxes")
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
