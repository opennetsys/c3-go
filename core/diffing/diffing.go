package diffing

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/c3systems/c3/common/command"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNoDifferencesFound ...
	ErrNoDifferencesFound = errors.New("no differences were found")

	// ErrCommandNotFound ...
	ErrCommandNotFound = errors.New("command not found")
)

// Diff ...
func Diff(old, new, out string, isDir bool) error {
	if !command.Exists("diff") {
		log.Println("[miner] error; diff command not found")
		return ErrCommandNotFound
	}

	var (
		commands []interface{}
		s        string
	)

	if isDir {
		commands = append(commands, "-rN")
		s = " %s"
	}

	commands = append(commands, old, new, ">", out)
	s += " %s %s %s %s"

	cmd := exec.Command("sh", "-c", fmt.Sprintf("diff -ud"+s, commands...))
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
	if !command.Exists("combinediff") {
		log.Println("[miner] error; combinediff command not found")
		return ErrCommandNotFound
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("combinediff %s %s > %s", firstDiff, secondDiff, out))
	if err := cmd.Start(); err != nil {
		log.Printf("[miner] error executing combinediff command; %s", err)
		return err
	}

	return cmd.Wait()
}

// Patch ...
func Patch(patch string, backup bool, absPath bool) error {
	if !command.Exists("patch") {
		log.Println("[miner] error; patch command not found")
		return ErrCommandNotFound
	}

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

	cmd := exec.Command("sh", "-c", fmt.Sprintf("patch"+s, commands...))
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
