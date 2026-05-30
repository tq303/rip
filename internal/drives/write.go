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
	bar := progress.Bar("b:", info.Size())

	destination, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer destination.Close()

	buffer := make([]byte, megaBytes*1024*1024)

	for {
		if state.Get().Err() != nil {
			fmt.Println("\nwarning: write was interrupted - drive maybe corrupted")
			os.Exit(1)
		}
		n, readErr := file.Read(buffer)
		if n > 0 {
			if _, writeErr := io.MultiWriter(destination, bar).Write(buffer[:n]); writeErr != nil {
				return writeErr
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	bar.Close()
	fmt.Print("\nWriting...")
	destination.Sync()
	fmt.Printf(" done in %s\n", time.Since(start).Round(time.Second))

	return nil
}
