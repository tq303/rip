package drives

import (
	"bytes"
	"encoding/xml"
	"os/exec"
	"strconv"
	"strings"
)

func listMac() ([]Drive, error) {
	list, err := exec.Command("diskutil", "list", "-plist", "external").Output()
	if err != nil {
		return nil, nil
	}
	drives, err := parsePlist(list)
	if err != nil {
		return nil, err
	}
	for i, d := range drives {
		info, err := exec.Command("diskutil", "info", "-plist", d.Name).Output()
		if err == nil {
			drives[i].Label = parseMediaName(info)
		}
	}
	return drives, nil
}

func UnmountMacOs(name string) error {
	return exec.Command("diskutil", "unmountDisk", name).Run()
}

func parseMediaName(data []byte) string {
	decode := xml.NewDecoder(bytes.NewReader(data))
	var currentElement, lastKey string
	for {
		token, err := decode.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
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

func parsePlist(data []byte) ([]Drive, error) {
	decode := xml.NewDecoder(bytes.NewReader(data))

	var result []Drive
	var currentKey, currentElement string
	var inTarget, inTopDict bool
	var current Drive
	depth, arrayDepth, dictDepth := 0, 0, 0

	for {
		token, err := decode.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
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
				if current.Name != "" && current.Size > 1024*1024 {
					result = append(result, current)
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
	return result, nil
}
