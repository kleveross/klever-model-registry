package util

import (
	"fmt"
	"os/exec"
)

func ExecOSCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	if b, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("exec %v %v is error: %v, %v", command, args, err, b)
	}
	return nil
}
