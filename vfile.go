package vhosts

import "os"

// utility functions

// doesFileExist checks if a file exists at the given path
func doesFileExist(path string) bool {
	// return true if the file already exists, if not return false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
