package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/schollz/progressbar/v3"
)

func downloadUrl(url string, outputFolder string) (string, error) {
	head, err := http.Head(url)
	if err != nil {
		// TODO check "Content-Type: application/x-iso9660-image"?
		return "", err
	}
	defer head.Body.Close()

	fileName := path.Base(head.Request.URL.Path)

	if outputFolder == "" {
		outputFolder = "/tmp/rip"
	}

	destination := path.Join(outputFolder, fileName)

	if out, err := os.Stat(destination); err == nil && out.Size() == head.ContentLength {
		fmt.Printf("Using cached %s\n", destination)
		return destination, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, url)
	}

	fmt.Printf("Downloading %s\n", fileName)

	if err := os.MkdirAll("/tmp/rip", 0755); err != nil {
		return "", err
	}
	out, err := os.Create(destination)
	if err != nil {
		return "", err
	}
	defer out.Close()

	progress := progressbar.New64(resp.ContentLength)

	_, err = io.Copy(out, io.TeeReader(resp.Body, progress))
	if err != nil {
		os.Remove(out.Name())
	}

	progress.Close()
	fmt.Printf("\nComplete %s\n", destination)

	return out.Name(), err
}

func getReleases(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, url)
	}

	// get all links
	// list all avaiable

	return "", nil
}
