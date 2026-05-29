package drives

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/schollz/progressbar/v3"
)

func Write(image string, target string, mb int) error {
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
	progress := progressbar.DefaultBytes(info.Size())

	destination, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer destination.Close()

	buffer := make([]byte, mb*1024*1024)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nwarning: write was interrupted - drive maybe corrupted")
		os.Exit(1)
	}()

	_, err = io.CopyBuffer(io.MultiWriter(destination, progress), file, buffer)
	if err != nil {
		return err
	}

	fmt.Printf("\nwritten in %s\n", time.Since(start).Round(time.Second))

	destination.Sync()

	return nil
}
