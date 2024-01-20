package main

import (
	"errors"
	"fmt"
	"os"
)

var CONF_DIR string

func init() {
	home, _ := os.UserHomeDir()

	CONF_DIR = fmt.Sprintf("%s/.utok", home)

	// Ignore errors: it just means the directory exists..
	os.Mkdir(CONF_DIR, 0755)
}

func writeFile(name string, data []byte) error {
	return os.WriteFile(fmt.Sprintf("%s/%s", CONF_DIR, name), data, 0600)
}

func readFile(name string) ([]byte, error) {
	return os.ReadFile(fmt.Sprintf("%s/%s", CONF_DIR, name))
}

func fileExists(name string) bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", CONF_DIR, name)); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func removeFile(name string) error {
	return os.Remove(fmt.Sprintf("%s/%s", CONF_DIR, name))
}
