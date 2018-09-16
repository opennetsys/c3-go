package diffing

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/c3systems/c3-go/common/command"
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
		log.Error("[miner] error; diff command not found")
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

	oldFile := filepath.Base(old)

	/*
	 * note: the first two lines of a standard patch file are:
	 *    --- oldFileName timestamp
	 *    +++ newFileName timestamp
	 *
	 *    We want the oldFileName to show up on both lines, so we replace the first instance in the file with sed
	 *
	 */
	commands = append(commands, old, new, old, oldFile, new, oldFile, out)
	s += " %s %s | sed -e 's|%s|%s|' | sed -e 's|%s|%s|' > %s"

	// same as: git diff --minimal --unified file1 file2
	cmd := exec.Command("sh", "-c", fmt.Sprintf("diff -uda"+s, commands...))
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

// CombineDiff ...
func CombineDiff(firstDiff, secondDiff, out string) error {
	if !command.Exists("combinediff") {
		log.Error("[miner] error; combinediff command not found")
		return ErrCommandNotFound
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("combinediff %s %s > %s", firstDiff, secondDiff, out))
	if err := cmd.Start(); err != nil {
		log.Errorf("[miner] error executing combinediff command; %s", err)
		return err
	}

	return cmd.Wait()
}

// Patch ...
func Patch(patch, orig string, backup bool, absPath bool) error {
	if !command.Exists("patch") {
		log.Error("[miner] error; patch command not found")
		return ErrCommandNotFound
	}

	var (
		commands []interface{}
		s        string
	)

	commands = append(commands, "-u")
	s = " %s"

	if backup {
		commands = append(commands, "-b")
		s += " %s"
	}
	if absPath {
		commands = append(commands, "-d/ -p0")
		s += " %s"
	}

	commands = append(commands, orig, patch)
	s += " %s %s"

	cmd := exec.Command("sh", "-c", fmt.Sprintf("patch"+s, commands...))
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
