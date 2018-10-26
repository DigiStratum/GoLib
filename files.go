package golib

/*

File handling library functions

*/

import(
	"fmt"
	"io/ioutil"
	"path/filepath"
	"errors"
)

func ReadFileString(path string) (*string, error) {
	tbuf, err := ReadFileBytes(path)
	if nil != err { return nil, err }
	s := string(*tbuf)
	return &s, nil
}

func ReadFileBytes(path string) (*[]byte, error) {
	var tbuf []byte
	tbuf, err = ioutil.ReadFile(path)
	if nil != err { return nil, errors.New(fmt.Sprintf("Error reading '%s'", path) }
	return &tbuf, nil
}

