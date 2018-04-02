package main

import (
	"log"
	"os"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var statusDetailTemplate = `
{{- range $i, $e := . -}}
--[ {{add $i 1}} ]-------------------------------------
       id: {{.ID}}
     name: {{.Name}}
 provider: {{.Provider}}
    state: {{.State}}
directory: {{.Directory}}
    vm_id: {{.VMID}}

{{end -}}
`

func statusCommand(c *cli.Context) {
	vms, err := GlobalStatus()
	if err != nil {
		log.Fatal(err)
	}
	if c.Bool("vertical") {
		printStatusVertically(vms)
	} else {
		printStatus(vms)
	}
}

func printStatus(vms []VagrantVM) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetHeaderLine(true)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{"id", "name", "provider", "state", "directory", "vm_id"})
	for _, vm := range vms {
		table.Append([]string{
			vm.ShortID,
			vm.Name,
			vm.Provider,
			vm.State,
			vm.Directory,
			vm.VMID,
		})
	}
	table.Render()
}

func printStatusVertically(vms []VagrantVM) {
	tmpl := template.New("status").Funcs(template.FuncMap{
		"add": tmplFuncAdd,
	})
	tmpl = template.Must(tmpl.Parse(statusDetailTemplate))
	err := tmpl.Execute(os.Stdout, vms)
	if err != nil {
		log.Fatal(err)
	}
}

func tmplFuncAdd(a, b int) int {
	return a + b
}
