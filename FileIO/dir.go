package fileio

import (
	"os"
	"io/fs"
	"path/filepath"
	"regexp"
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
	var patternRexp *regexp.Regexp
	var err error
	if len(pattern) > 0 {
		patternRexp, err = regexp.Compile(pattern)
		if nil != err { return nil, err }
	}
	files := NewFileSet()
	if err = filepath.Walk(
		r.path,
		func (file string, f os.FileInfo, err error) error {
			// Fail on error
			if nil != err { return err }
			// No match on non-files
			if ! f.IsDir() { return nil }
			// No match on pattern regex if specified
			if (nil != patternRexp) && (! patternRexp.MatchString(file)) { return nil }
			// Add to matches
			files.AddFile(file)
			return nil
		},
	); nil != err { return nil, err }
	return files, nil
}

// -------------------------------------------------------------------------------------------------
// dir
// -------------------------------------------------------------------------------------------------

// Get the FileInfo for this File with a read-through local cache copy
func (r *dir) getFileInfo() (*fs.FileInfo, error) {
	// If we don't have a cached copy already
	if (nil == r.fileInfo) {
		// Pull FileInfo from the os
		var err error
		var fi fs.FileInfo
		if fi, err = os.Stat(r.path); nil != err { return nil, err }
		// And cache the result
		r.fileInfo = &fi
	}
	// Return the cached result
	return r.fileInfo, nil
}


