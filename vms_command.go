package main

import (
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func vmsCommand(c *cli.Context) {
	vms, err := GetVirtualBoxVMs()
	if err != nil {
		log.Fatal(err)
	}
	printVMs(vms)
}

func printVMs(vms []VirtualBoxVM) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetHeaderLine(true)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{"vm_id", "name"})

	for _, vm := range vms {
		table.Append([]string{
			vm.ID,
			vm.Name,
		})
	}
	table.Render()
}
