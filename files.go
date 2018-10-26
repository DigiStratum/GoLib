package golib

/*

File handling library functions

*/

import(
	"os"
	"fmt"
	"io/ioutil"
	"errors"
)

// Write the contents of a string to a file
func WriteFileString(path string, content *string) error {
	c := []byte(*content)
	return ioutil.WriteFile(path, c, 0644)
}

// Read the file located at the specified path and return the contents as a *string
func ReadFileString(path string) (*string, error) {
	tbuf, err := ReadFileBytes(path)
	if nil != err { return nil, err }
	s := string(*tbuf)
	return &s, nil
}

// Read the file located at the specified path and return the contents as a *[]byte
func ReadFileBytes(path string) (*[]byte, error) {
	var tbuf []byte
	var err error
	tbuf, err = ioutil.ReadFile(path)
	if nil != err {
		return nil, errors.New(fmt.Sprintf("Error reading '%s': %s", path, err.Error()))
	}
	return &tbuf, nil
}

// Is the specified path a directory? return true if so, else false
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if nil != err { return false } // Who knows?
	mode := fi.Mode();
	return mode.IsDir()
}

// Is the specified path a false? return true if so, else false
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if nil != err { return false } // Who knows?
	mode := fi.Mode();
	return mode.IsRegular()
}

