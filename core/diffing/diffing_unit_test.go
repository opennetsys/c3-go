// +build unit

package diffing

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	// TODO: test dir
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("err getting current working directory\n%v", err)
	}

	if err := createDirIfNotExist(fmt.Sprintf("%s/.testfiles/", wd)); err != nil {
		t.Errorf("err creating test files directory\n%v", err)
	}
	//defer os.RemoveAll(fmt.Sprintf("%s/.testfiles", wd))

	f1, err := os.Create(fmt.Sprintf("%s/.testfiles/1.txt", wd))
	if err != nil {
		t.Errorf("err creating first file\n%v", err)
	}
	defer f1.Close()

	if err := f1.Chmod(0777); err != nil {
		t.Errorf("eff chmodding file 1\n%v", err)
	}

	if _, err := f1.WriteString("1"); err != nil {
		t.Errorf("err writing to first file\n%v", err)
	}

	f1.Sync()

	f2, err := os.Create(fmt.Sprintf("%s/.testfiles/2.txt", wd))
	if err != nil {
		t.Errorf("err creating second file\n%v", err)
	}
	defer f2.Close()

	if err := f2.Chmod(0777); err != nil {
		t.Errorf("eff chmodding file 2\n%v", err)
	}

	if _, err := f2.WriteString("12"); err != nil {
		t.Errorf("err writing to second file\n%v", err)
	}

	f2.Sync()

	if err := Diff(fmt.Sprintf("%s/.testfiles/1.txt", wd), fmt.Sprintf("%s/.testfiles/2.txt", wd), fmt.Sprintf("%s/.testfiles/12.patch", wd), false); err != nil {
		t.Errorf("err diffing the files\n%v", err)
	}

	actual, err := ioutil.ReadFile(fmt.Sprintf("%s/.testfiles/12.patch", wd))
	if err != nil {
		t.Errorf("err reading patch file\n%v", err)
	}

	expected := `@@ -1 +1 @@
-1
\ No newline at end of file
+12
\ No newline at end of file`

	if !strings.Contains(string(actual), expected) {
		t.Errorf("expected\n%s\nin\n%s", expected, string(actual))
	}
}

func TestCombineDiff(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("err getting current working directory\n%v", err)
	}

	if err := createDirIfNotExist(fmt.Sprintf("%s/.testfiles/", wd)); err != nil {
		t.Errorf("err creating test files directory\n%v", err)
	}
	defer os.RemoveAll(fmt.Sprintf("%s/.testfiles", wd))

	f1, err := os.Create(fmt.Sprintf("%s/.testfiles/1.txt", wd))
	if err != nil {
		t.Errorf("err creating first file\n%v", err)
	}
	defer f1.Close()

	if err := f1.Chmod(0777); err != nil {
		t.Errorf("eff chmodding file 1\n%v", err)
	}

	if _, err := f1.WriteString("1"); err != nil {
		t.Errorf("err writing to first file\n%v", err)
	}

	if err := f1.Sync(); err != nil {
		t.Errorf("err synching file 1\n%v", err)
	}

	f2, err := os.Create(fmt.Sprintf("%s/.testfiles/2.txt", wd))
	if err != nil {
		t.Errorf("err creating second file\n%v", err)
	}
	defer f2.Close()

	if err := f2.Chmod(0777); err != nil {
		t.Errorf("eff chmodding file 2\n%v", err)
	}

	if _, err := f2.WriteString("12"); err != nil {
		t.Errorf("err writing to second file\n%v", err)
	}

	if err := f2.Sync(); err != nil {
		t.Errorf("err synching file 2\n%v", err)
	}

	f3, err := os.Create(fmt.Sprintf("%s/.testfiles/3.txt", wd))
	if err != nil {
		t.Errorf("err creating second file\n%v", err)
	}
	defer f3.Close()

	if err := f3.Chmod(0777); err != nil {
		t.Errorf("eff chmodding file 3\n%v", err)
	}

	if _, err := f3.WriteString("123"); err != nil {
		t.Errorf("err writing to third file\n%v", err)
	}

	if err := f3.Sync(); err != nil {
		t.Errorf("err synching file 3\n%v", err)
	}

	if err := Diff(fmt.Sprintf("%s/.testfiles/1.txt", wd), fmt.Sprintf("%s/.testfiles/2.txt", wd), fmt.Sprintf("%s/.testfiles/12.patch", wd), false); err != nil {
		t.Errorf("err diffing the 1 and 2\n%v", err)
	}

	if err := Diff(fmt.Sprintf("%s/.testfiles/2.txt", wd), fmt.Sprintf("%s/.testfiles/3.txt", wd), fmt.Sprintf("%s/.testfiles/23.patch", wd), false); err != nil {
		t.Errorf("err diffing the 2 and 3\n%v", err)
	}

	if err := CombineDiff(fmt.Sprintf("%s/.testfiles/12.patch", wd), fmt.Sprintf("%s/.testfiles/23.patch", wd), fmt.Sprintf("%s/.testfiles/13.combined.patch", wd)); err != nil {
		t.Errorf("err combining diffs\n%v", err)
	}

	actual, err := ioutil.ReadFile(fmt.Sprintf("%s/.testfiles/13.combined.patch", wd))
	if err != nil {
		t.Errorf("err reading patch file\n%v", err)
	}

	expected1 := `@@ -1 +1 @@
-1
\ No newline at end of file
+12`
	expected2 := `@@ -1 +1 @@
-12
\ No newline at end of file
+123`

	if !strings.Contains(string(actual), expected1) ||
		!strings.Contains(string(actual), expected2) {
		t.Errorf("expected\n%s\nand\n%s\nin\n%s", expected1, expected2, string(actual))
	}
}

func TestPatch(t *testing.T) {
	// TODO: add tests for not backup and not absPath
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("err getting current working directory\n%v", err)
	}

	if err := createDirIfNotExist(fmt.Sprintf("%s/.testfiles/", wd)); err != nil {
		t.Errorf("err creating test files directory\n%v", err)
	}
	defer os.RemoveAll(fmt.Sprintf("%s/.testfiles", wd))

	f1, err := os.Create(fmt.Sprintf("%s/.testfiles/1.txt", wd))
	if err != nil {
		t.Errorf("err creating first file\n%v", err)
	}
	defer f1.Close()

	if err := f1.Chmod(0777); err != nil {
		t.Errorf("eff chmodding file 1\n%v", err)
	}

	if _, err := f1.WriteString("1"); err != nil {
		t.Errorf("err writing to first file\n%v", err)
	}

	f1.Sync()

	f2, err := os.Create(fmt.Sprintf("%s/.testfiles/2.txt", wd))
	if err != nil {
		t.Errorf("err creating second file\n%v", err)
	}
	defer f2.Close()

	if err := f2.Chmod(0777); err != nil {
		t.Errorf("eff chmodding file 2\n%v", err)
	}

	if _, err := f2.WriteString("12"); err != nil {
		t.Errorf("err writing to second file\n%v", err)
	}

	f2.Sync()

	if err := Diff(fmt.Sprintf("%s/.testfiles/1.txt", wd), fmt.Sprintf("%s/.testfiles/2.txt", wd), fmt.Sprintf("%s/.testfiles/12.patch", wd), false); err != nil {
		t.Errorf("err diffing the files\n%v", err)
	}

	if err := Patch(fmt.Sprintf("%s/.testfiles/12.patch", wd), fmt.Sprintf("%s/.testfiles/1.txt", wd), true, true); err != nil {
		t.Errorf("err patching\n%v", err)
	}

	actual, err := ioutil.ReadFile(fmt.Sprintf("%s/.testfiles/1.txt", wd))
	if err != nil {
		t.Errorf("err reading patch file\n%v", err)
	}
	backup, err := ioutil.ReadFile(fmt.Sprintf("%s/.testfiles/1.txt.orig", wd))
	if err != nil {
		t.Errorf("err reading patch file\n%v", err)
	}

	expected := "12"
	expectedBackup := "1"

	if expected != string(actual) || string(backup) != expectedBackup {
		t.Errorf("expected\n%s\nand\n%s\nreceived\n%s\nand\n%s", expected, expectedBackup, string(actual), string(backup))
	}
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}

	return nil
}
