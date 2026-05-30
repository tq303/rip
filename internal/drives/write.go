package drives

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tq303/rip/internal/progress"
	"github.com/tq303/rip/internal/state"
)

func Write(image string, target string, megaBytes int) error {
	info, err := os.Stat(image)
	if err != nil {
		return err
	}

	file, err := os.Open(image)
	if err != nil {
		return err
	}
	defer file.Close()

	start := time.Now()
	bar := progress.Bar("buffering", info.Size())

	destination, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer destination.Close()

	buffer := make([]byte, megaBytes*1024*1024)

	if state.Get().Err() != nil {
		fmt.Println("\nwarning: write was interrupted - drive maybe corrupted")
		os.Exit(1)
	}

	_, err = io.CopyBuffer(io.MultiWriter(destination, bar), file, buffer)
	if err != nil {
		return err
	}

	bar.Close()
	fmt.Print("\nWriting...")
	destination.Sync()
	fmt.Printf(" done in %s\n", time.Since(start).Round(time.Second))

	return nil
}
