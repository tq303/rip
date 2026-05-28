package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Drive struct {
	Name  string
	Path  string
	Size  uint64
	Label string
}

func listDrives() ([]Drive, error) {
	switch runtime.GOOS {
	case "darwin":
		return listDrivesMac()
	case "linux":
		return listDrivesLinux()
	case "windows":
		return listDrivesWindows()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func listDrivesMac() ([]Drive, error) {
	out, err := exec.Command("diskutil", "list", "-plist", "external").Output()
	if err != nil {
		return nil, nil
	}
	drives, err := parseMacPlist(out)
	if err != nil {
		return nil, err
	}
	for i, d := range drives {
		info, err := exec.Command("diskutil", "info", "-plist", d.Name).Output()
		if err == nil {
			drives[i].Label = parseMacMediaName(info)
		}
	}
	return drives, nil
}

func parseMacMediaName(data []byte) string {
	dec := xml.NewDecoder(bytes.NewReader(data))
	var currentElement, lastKey string
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			currentElement = t.Name.Local
		case xml.EndElement:
			currentElement = ""
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text == "" {
				continue
			}
			if currentElement == "key" {
				lastKey = text
			} else if currentElement == "string" && lastKey == "MediaName" {
				return text
			}
		}
	}
	return ""
}

func parseMacPlist(data []byte) ([]Drive, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))

	var drives []Drive
	var currentKey, currentElement string
	var inTarget, inTopDict bool
	var current Drive
	depth, arrayDepth, dictDepth := 0, 0, 0

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			currentElement = t.Name.Local
			if t.Name.Local == "array" && currentKey == "AllDisksAndPartitions" {
				inTarget = true
				arrayDepth = depth
			}
			if inTarget && !inTopDict && t.Name.Local == "dict" && depth == arrayDepth+1 {
				inTopDict = true
				dictDepth = depth
				current = Drive{}
			}
		case xml.EndElement:
			if inTopDict && t.Name.Local == "dict" && depth == dictDepth {
				if current.Name != "" {
					drives = append(drives, current)
				}
				inTopDict = false
			}
			if inTarget && t.Name.Local == "array" && depth == arrayDepth {
				inTarget = false
			}
			depth--
			currentElement = ""
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text == "" {
				continue
			}
			if currentElement == "key" {
				currentKey = text
				continue
			}
			if inTopDict && depth == dictDepth+1 {
				switch currentKey {
				case "DeviceIdentifier":
					if currentElement == "string" {
						current.Name = text
						current.Path = "/dev/r" + text
						currentKey = ""
					}
				case "Size":
					if currentElement == "integer" {
						size, _ := strconv.ParseUint(text, 10, 64)
						current.Size = size
						currentKey = ""
					}
				}
			}
		}
	}
	return drives, nil
}

func unmountDisk(name string) error {
	return exec.Command("diskutil", "unmountDisk", name).Run()
}

func listDrivesLinux() ([]Drive, error) {
	_, err := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,TRAN,RM").Output()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func listDrivesWindows() ([]Drive, error) {
	_, err := exec.Command("powershell", "-Command", "Get-Disk | Select-Object Number,FriendlyName,Size,BusType | ConvertTo-Json").Output()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
