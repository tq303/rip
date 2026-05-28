package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
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

	options := make([]huh.Option[string], len(drives))
	for i, d := range drives {
		label := fmt.Sprintf("%s — %s — %.1f GB", d.Label, d.Path, float64(d.Size)/1e9)
		options[i] = huh.NewOption(label, d.Path)
	}

	var target string
	err = huh.NewSelect[string]().
		Title("Select a drive").
		Options(options...).
		Value(&target).
		Run()
	if err != nil {
		return err
	}

	fmt.Println("target:", target)
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
