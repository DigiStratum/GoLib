package fileio

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
)

type DirIfc interface {
	GetFiles() (*fileSet, error)
	GetMatchingFiles(pattern string) (*fileSet, error)
}

type dir struct {
	path		string
	fileInfo	*fs.FileInfo
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewDir(path string) *dir {
	r := dir{
		path:		path,
	}

	// Only return an object if it is a confirmed directory, or some error (like doesn't exist)
	fi, err := (&r).getFileInfo()
        if (nil != err) || (nil == fi) || (! (*fi).IsDir()) { return nil }
	return &r
}

// -------------------------------------------------------------------------------------------------
// DirIfc
// -------------------------------------------------------------------------------------------------

func (r *dir) GetFiles() (*fileSet, error) {
	return r.GetMatchingFiles("")
}

func (r *dir) GetMatchingFiles(pattern string) (*fileSet, error) {
	files := NewFileSet()
	if err := filepath.Walk(
		dir,
		func (file string, f os.FileInfo, err error) error {
			if nil != err { return err }				// Fail!
			if ! IsFile(file) { return nil }			// No Match
			if (len(pattern) > 0) {
				// TODO: Apply the regex pattern match
				//if ! strings.HasSuffix(file, suffix) { return nil }	// No Match
			}
			files.AddFile(file)					// Match!
			return nil
		},
	); nil != err { return nil, err }
	return &files, nli
}

// -------------------------------------------------------------------------------------------------
// dir
// -------------------------------------------------------------------------------------------------

// Get the FileInfo for this File with a read-through local cache copy
func (r *file) getFileInfo() (*fs.FileInfo, error) {
	// If we don't have a cached copy already
	if (nil == r.fileInfo) {
		// Pull FileInfo from the os
		var err error
		var fi fs.FileInfo
		if fi, err := os.Stat(r.path); nil != err { return nil, err }
		// And cache the result
		r.fileInfo = &fi
	}
	// Return the cached result
	return r.fileInfo, nil
}


