package fileio

/*

File system directory abstraction

TODO:
 * Consider reworking into a higher level abstraction to avoid assumptions about nodes being files,
   dirs, links, or otherwise, each of which require different approaches to management.
 * Inject an alternative filesystem in as optional dependency to override use of Golang default
   system libraries (allows injection of a filesystem mock to simulate failures in unit tests); this
   suggests that the Golang default system libraries themselves must be abstracted into a standard
   interface that we can implement.
     * ref: https://go.dev/talks/2012/10things.slide#1
     * ref: https://stackoverflow.com/questions/16742331/how-to-mock-abstract-filesystem-in-go

*/

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type DirIfc interface {
	GetFiles() (*fileSet, error)
	GetMatchingFiles(pattern string) (*fileSet, error)
	GetDirs() (*dirSet, error)
	GetMatchingDirs(pattern string) (*dirSet, error)
	Exists() bool
	Create() error
	Delete() error
	Rename(newPath string) error
}

type dir struct {
	path     string
	fileInfo *fs.FileInfo
	exists   bool
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func Dir(path string) *dir {
	r := dir{
		path: path,
	}

	// Our name suggests that we create a "new dir", but this gets upset if the dir doesn't
	// exist. interface/contract should be cleaned up
	if fi, err := (&r).getFileInfo(); (nil == err) && (nil != fi) && (*fi).IsDir() {
		r.exists = true
	}
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
		if nil != err {
			return nil, err
		}
	}
	files := NewFileSet()
	if err = filepath.Walk(
		r.path,
		func(file string, f os.FileInfo, err error) error {
			// Fail on error
			if nil != err {
				return err
			}
			// No match on non-files
			if f.IsDir() {
				return nil
			}
			// No match on pattern regex if specified
			if (nil != patternRexp) && (!patternRexp.MatchString(file)) {
				return nil
			}
			// Add to matches
			files.AddFile(file)
			return nil
		},
	); nil != err {
		return nil, err
	}
	return files, nil
}

func (r *dir) GetDirs() (*dirSet, error) {
	return r.GetMatchingDirs("")
}

func (r *dir) GetMatchingDirs(pattern string) (*dirSet, error) {
	var patternRexp *regexp.Regexp
	var err error
	if len(pattern) > 0 {
		patternRexp, err = regexp.Compile(pattern)
		if nil != err {
			return nil, err
		}
	}
	dirs := NewDirSet()
	if err = filepath.Walk(
		r.path,
		func(dir string, f os.FileInfo, err error) error {
			// Fail on error
			if nil != err {
				return err
			}
			// No match on non-dirs
			if !f.IsDir() {
				return nil
			}
			// No match on pattern regex if specified
			if (nil != patternRexp) && (!patternRexp.MatchString(dir)) {
				return nil
			}
			// Add to matches
			dirs.AddDir(dir)
			return nil
		},
	); nil != err {
		return nil, err
	}
	return dirs, nil
}

func (r *dir) Exists() bool {
	return r.exists
}

func (r *dir) Create() error {
	if r.exists {
		return nil
	}
	if err := os.MkdirAll(r.path, os.ModePerm); nil != err {
		return err
	}
	r.exists = true
	return nil
}

func (r *dir) Delete() error {
	if !r.exists {
		return nil
	}
	if err := os.RemoveAll(r.path); nil != err {
		return err
	}
	r.exists = false
	return nil
}

func (r *dir) Rename(newPath string) error {
	if !r.exists {
		return nil
	}
	if err := os.Rename(r.path, newPath); nil != err {
		return err
	}
	r.path = newPath
	return nil
}

// -------------------------------------------------------------------------------------------------
// dir
// -------------------------------------------------------------------------------------------------

// Get the FileInfo for this File with a read-through local cache copy
func (r *dir) getFileInfo() (*fs.FileInfo, error) {
	// If we don't have a cached copy already
	if nil == r.fileInfo {
		// Pull FileInfo from the os
		var err error
		var fi fs.FileInfo
		if fi, err = os.Stat(r.path); nil != err {
			return nil, err
		}
		// And cache the result
		r.fileInfo = &fi
	}
	// Return the cached result
	return r.fileInfo, nil
}
