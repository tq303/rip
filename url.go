package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

func changeFileOwner(fileName string) {
	if uid := os.Getenv("SUDO_UID"); uid != "" {
		if gid := os.Getenv("SUDO_GID"); gid != "" {
			u, _ := strconv.Atoi(uid)
			g, _ := strconv.Atoi(gid)
			os.Chown(fileName, u, g)
		}
	}
}

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

	if localFile, err := os.Stat(destination); err == nil && localFile.Size() == head.ContentLength {
		fmt.Printf("Using cached %s\n\n", destination)
		return destination, nil
	}

	start := time.Now()
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
	saveFile, err := os.Create(destination)
	if err != nil {
		return "", err
	}
	defer saveFile.Close()

	progress := progressbar.DefaultBytes(resp.ContentLength)

	_, err = io.Copy(saveFile, io.TeeReader(resp.Body, progress))
	if err != nil {
		os.Remove(saveFile.Name())
	}

	progress.Close()

	fmt.Printf("Downloaded in %s\n\n", time.Since(start).Round(time.Second))

	changeFileOwner(saveFile.Name())

	return saveFile.Name(), err
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
