package drives

import (
	"encoding/json"
	"os/exec"
)

type linuxDevice struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type string `json:"type"`
	Tran string `json:"tran"`
	RM   bool   `json:"rm"`
}

type lsblk struct {
	BlockDevices []linuxDevice `json:"blockdevices"`
}

func listLinux() ([]Drive, error) {
	jsonOutput, err := exec.Command("lsblk", "-J", "-b", "-o", "NAME,SIZE,TYPE,TRAN,RM").Output()
	if err != nil {
		return nil, err
	}

	var devices lsblk
	err = json.Unmarshal(jsonOutput, &devices)
	if err != nil {
		return nil, err
	}

	options := []Drive{}

	for _, device := range devices.BlockDevices {
		if device.Type == "disk" && device.RM {
			options = append(options, Drive{
				Name: device.Name,
				Path: "/dev/" + device.Name,
				Size: device.Size,
			})
		}
	}

	return options, nil
}
