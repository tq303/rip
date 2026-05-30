package drives

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type winDisk struct {
	Number       int    `json:"Number"`
	FriendlyName string `json:"FriendlyName"`
	Size         uint64 `json:"Size"`
	BusType      string `json:"BusType"`
}

func listWindows() ([]Drive, error) {
	out, err := exec.Command("powershell", "-Command", "Get-Disk | Select-Object Number,FriendlyName,Size,BusType | ConvertTo-Json").Output()
	if err != nil {
		return nil, err
	}

	// PowerShell returns a single object (not array) when only one disk is present
	var disks []winDisk
	if err := json.Unmarshal(out, &disks); err != nil {
		var single winDisk
		if err := json.Unmarshal(out, &single); err != nil {
			return nil, fmt.Errorf("parsing disk list: %w", err)
		}
		disks = []winDisk{single}
	}

	var result []Drive
	for _, d := range disks {
		if d.BusType != "USB" && d.BusType != "SD" && d.BusType != "MMC" {
			continue
		}
		result = append(result, Drive{
			Name:  fmt.Sprintf("PhysicalDrive%d", d.Number),
			Path:  fmt.Sprintf(`\\.\PhysicalDrive%d`, d.Number),
			Size:  d.Size,
			Label: d.FriendlyName,
		})
	}
	return result, nil
}
