package drives

import "os/exec"

func listWindows() ([]Drive, error) {
	_, err := exec.Command("powershell", "-Command", "Get-Disk | Select-Object Number,FriendlyName,Size,BusType | ConvertTo-Json").Output()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
