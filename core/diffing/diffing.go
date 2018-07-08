package diffing

import (
	"errors"
	"fmt"
	"os/exec"
)

var (
	// ErrNoDifferencesFound ...
	ErrNoDifferencesFound = errors.New("no differences were found")
)

// Diff ...
func Diff(old, new, out string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("/usr/bin/diff -ud %s %s > %s", old, new, out))
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		/*
		 * note:
		 *   0 == no diffs found;
		 *   1 == diffs were found;
		 *   2 == some other err;
		 *
		 */
		if err.Error() == "exit status 1" {
			return nil
		}

		return err
	}

	return ErrNoDifferencesFound
}

// CombineDiff ...
func CombineDiff(firstDiff, secondDiff, out string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("/usr/bin/combinediff %s %s > %s", firstDiff, secondDiff, out))
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

// Patch ...
func Patch(patch string, backup bool, absPath bool) error {
	var (
		commands []interface{}
		s        string
	)

	if backup {
		commands = append(commands, "-b")
		s = " %s"
	}
	if absPath {
		commands = append(commands, "-d/ -p0")
		s += " %s"
	}

	commands = append(commands, "<", patch)
	s += " %s %s"

	cmd := exec.Command("sh", "-c", fmt.Sprintf("/usr/bin/patch"+s, commands...))
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
