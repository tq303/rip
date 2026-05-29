package drives

import (
	"fmt"
	"runtime"
)

type Drive struct {
	Name  string
	Path  string
	Size  uint64
	Label string
}

func List() ([]Drive, error) {
	switch runtime.GOOS {
	case "darwin":
		return listMac()
	case "linux":
		return listLinux()
	case "windows":
		return listWindows()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
