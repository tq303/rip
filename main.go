package main

import (
	"fmt"
	"runtime"

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
	image := args[0]

	drives, err := listDrives()
	if err != nil {
		return err
	}
	if len(drives) == 0 {
		fmt.Println("no drives available")
		return nil
	}

	options := make([]huh.Option[Drive], len(drives))
	for i, d := range drives {
		label := fmt.Sprintf("%s — %s — %.1f GB", d.Label, d.Path, float64(d.Size)/1e9)
		options[i] = huh.NewOption(label, d)
	}

	var target Drive
	err = huh.NewSelect[Drive]().
		Title("Select a drive").
		Options(options...).
		Value(&target).
		Run()
	if err != nil {
		return err
	}

	var confirm bool
	err = huh.NewConfirm().
		Title(fmt.Sprintf("Flash %s to %s?", image, target.Label)).
		Description("This will erase all data on the drive.").
		Value(&confirm).
		Run()
	if err != nil {
		return err
	}
	if !confirm {
		return nil
	}

	if runtime.GOOS == "darwin" {
		if err := unmountDisk(target.Name); err != nil {
			return err
		}
	}

	return Write(image, target.Path)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
