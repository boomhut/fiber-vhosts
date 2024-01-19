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

// openFileToRead opens a file at the given path for reading. If the file doesn't exist, it returns an error
func openFileToRead(path string) (*os.File, error) {

	// Check if the file already exists, if it doesn't return an error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	// Open the file at the given path
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil

}
