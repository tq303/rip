//go:build windows

package drives

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
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

const (
	fsctlLockVolume     = 0x00090018
	fsctlDismountVolume = 0x00090020
)

func UnmountWindows(path string) error {
	numStr := strings.TrimPrefix(path, `\\.\PhysicalDrive`)
	out, err := exec.Command("powershell", "-Command",
		fmt.Sprintf("Get-Disk -Number %s | Get-Partition | Select-Object -ExpandProperty DriveLetter", numStr),
	).Output()
	if err != nil {
		return nil
	}

	for _, letter := range strings.Fields(string(out)) {
		volumePath := `\\.\` + letter + `:`
		h, err := windows.CreateFile(
			windows.StringToUTF16Ptr(volumePath),
			windows.GENERIC_READ|windows.GENERIC_WRITE,
			windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
			nil,
			windows.OPEN_EXISTING,
			0,
			0,
		)
		if err != nil {
			continue
		}
		var n uint32
		windows.DeviceIoControl(h, fsctlLockVolume, nil, 0, nil, 0, &n, nil)
		windows.DeviceIoControl(h, fsctlDismountVolume, nil, 0, nil, 0, &n, nil)
		// keep handle open until process exits — closing releases the lock
		_ = h
	}

	return nil
}
