//go:build !windows

package drives

import "fmt"

func listWindows() ([]Drive, error) {
	return nil, fmt.Errorf("unsupported platform")
}

func UnmountWindows(path string) error {
	return fmt.Errorf("unsupported platform")
}
