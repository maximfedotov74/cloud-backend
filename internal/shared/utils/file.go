package utils

import (
	"path/filepath"
	"strings"
)

func GetFileExt(filename string) string {
	return strings.TrimPrefix(filepath.Ext(filename), filepath.Base(filename))
}

func GetFileName(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}
