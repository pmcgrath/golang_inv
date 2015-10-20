package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetAllSubDirectoryPaths(t *testing.T) {
	tempDirPath := path.Join(os.TempDir(), fmt.Sprintf("fs_tests_%d", os.Getpid()))
	emptyDirPath := path.Join(tempDirPath, "empty")
	dirWithContentPath := path.Join(tempDirPath, "content")
	subDirs := []string{path.Join(dirWithContentPath, "dir1"), path.Join(dirWithContentPath, "dir2")}

	for _, dirPath := range []string{emptyDirPath, subDirs[0], subDirs[1]} {
		os.MkdirAll(dirPath, 0777)
	}
	os.Create(path.Join(dirWithContentPath, "afile"))
	defer os.RemoveAll(tempDirPath)

	for _, testCase := range []struct {
		Path            string
		ExpectedSubDirs []string
	}{
		{emptyDirPath, []string{}},
		{dirWithContentPath, subDirs},
	} {
		actual, err := getAllSubDirectoryPaths(testCase.Path)
		if err != nil {
			t.Errorf("Unexpected error encountered for [%s], got an error = %s", testCase.Path, err)
		}
		if !reflect.DeepEqual(actual, testCase.ExpectedSubDirs) {
			t.Errorf("Unexpected subdirs encountered for [%s], got %v but expected %v", testCase.Path, actual, testCase.ExpectedSubDirs)
		}
	}
}

func TestTestIfDirectoryExists(t *testing.T) {
	programFilePath, _ := filepath.Abs(os.Args[0])
	programDirPath := filepath.Dir(programFilePath)

	if testIfDirectoryExists(programFilePath) {
		t.Errorf("Did not expect true for [%s]", programFilePath)
	}

	if !testIfDirectoryExists(programDirPath) {
		t.Errorf("Did not expect false for [%s]", programDirPath)
	}
}
