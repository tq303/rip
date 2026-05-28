package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rip [image]",
	Short: "Flash an OS image to a drive",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	drives, err := listDrives()
	if err != nil {
		return err
	}
	if len(drives) == 0 {
		fmt.Println("no drives available")
		return nil
	}
	for _, d := range drives {
		fmt.Printf("%s  %s  %d bytes\n", d.Path, d.Label, d.Size)
	}
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
