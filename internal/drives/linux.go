package drives

import "os/exec"

func listLinux() ([]Drive, error) {
	_, err := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,TRAN,RM").Output()
	if err != nil {
		return nil, err
	}
	return nil, nil
}
