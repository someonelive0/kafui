package backend

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetPrgDir() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret, nil
}

func Chdir2PrgPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	os.Chdir(ret)
	return ret, nil
}
