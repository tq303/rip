package main

import (
	"fmt"
	"io"
	"os"
)

func Write(image string, target string) error {
	file, err := os.Open(image)
	if err != nil {
		return err
	}
	defer file.Close()

	destination, err := os.OpenFile(target, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer destination.Close()

	buffer := make([]byte, 1024*1024)

	_, err = io.CopyBuffer(destination, file, buffer)
	if err != nil {
		return err
	}

	fmt.Printf("image %s written to %s", image, target)

	destination.Sync()

	return nil
}
