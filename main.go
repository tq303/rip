package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

type Drive struct {
	Name  string
	Path  string
	Size  uint64
	Label string
}

var rootCmd = &cobra.Command{
	Use:   "rip [image]",
	Short: "Flash an OS image to a drive",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Println("platform:", runtime.GOOS)
	return nil
}

func main() {
	rootCmd.Execute()
}
