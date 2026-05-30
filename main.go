package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/tq303/rip/internal/download"
	"github.com/tq303/rip/internal/drives"
	"github.com/tq303/rip/internal/permissions"
)

var rootCmd = &cobra.Command{
	Use:   "rip [image] [flags]",
	Short: "Flash an image to a drive",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	start := time.Now()
	image := args[0]

	if strings.HasPrefix(image, "http") {
		outputFolder, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		tempFile, err := download.Url(image, outputFolder)
		if err != nil {
			return err
		}

		image = tempFile
	}

	info, err := os.Stat(image)
	if err != nil {
		return err
	}

	list, err := drives.List()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Println("no drives available")
		return nil
	}

	options := make([]huh.Option[drives.Drive], len(list))
	for i, d := range list {
		label := fmt.Sprintf("%s — %s — %.1f GB", d.Label, d.Path, float64(d.Size)/1e9)
		options[i] = huh.NewOption(label, d)
	}

	var target drives.Drive
	err = huh.NewSelect[drives.Drive]().
		Title("Select a drive").
		Options(options...).
		Value(&target).
		Run()
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	var confirm bool

	if !force {
		err = huh.NewConfirm().
			Title(fmt.Sprintf("Flash %s (%.1f GB) to %s (%.1f GB)?", image, float64(info.Size())/1e9, target.Label, float64(target.Size)/1e9)).
			Description("This will erase all data on the drive.").
			Value(&confirm).
			Run()
		if err != nil {
			return err
		}
	} else {
		confirm = force
	}

	if !confirm {
		return nil
	}

	if runtime.GOOS == "darwin" {
		if err := drives.UnmountMacOs(target.Name); err != nil {
			return err
		}
	}

	buffer, err := cmd.Flags().GetInt("buffer")
	if err != nil {
		return err
	}

	if err := drives.Write(image, target.Path, buffer); err != nil {
		return err
	}
	fmt.Printf("Total time: %s\n", time.Since(start).Round(time.Second))
	return nil
}

func main() {
	permissions.EnsureRoot()
	rootCmd.Flags().IntP("buffer", "b", 4, "Set write buffer size in MB")
	rootCmd.Flags().StringP("output", "o", "", "Set download folder for URLs")
	rootCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
