package permissions

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func EnsureRoot() {
	if os.Getuid() == 0 {
		return
	}
	if runtime.GOOS == "windows" {
		fmt.Fprintln(os.Stderr, "please re-run as administrator")
		os.Exit(1)
	}
	bin, err := exec.LookPath(os.Args[0])
	if err != nil {
		bin = os.Args[0]
	}
	cmd := exec.Command("sudo", append([]string{bin}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
