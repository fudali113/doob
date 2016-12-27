package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}

func IsDirectory(path string) bool {
	return strings.HasSuffix(path, "/")
}
