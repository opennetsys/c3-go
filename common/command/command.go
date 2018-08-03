package command

import "os/exec"

// Exists check if command exists
func Exists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}

	return true
}
