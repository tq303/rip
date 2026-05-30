package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/tq303/rip/internal/progress"
	"github.com/tq303/rip/internal/state"
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
	headReq, err := http.NewRequestWithContext(state.Get(), "HEAD", url, nil)
	if err != nil {
		return "", err
	}
	head, err := http.DefaultClient.Do(headReq)
	if err != nil {
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
	req, err := http.NewRequestWithContext(state.Get(), "GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
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

	bar := progress.Bar("downloading", resp.ContentLength)

	_, err = io.Copy(saveFile, io.TeeReader(resp.Body, bar))
	if err != nil {
		os.Remove(saveFile.Name())
	}

	bar.Close()

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
