package fileio

import (
	"os"
	"io/fs"
	"time"
)

type FileIfc interface {
	// FileInfo properties
	GetName() (*string, error)
	GetSize() (*int64, error)
	GetMode() (*fs.FileMode, error)
	GetModTime() (*time.Time, error)
	IsDir() (bool, error)
	GetSys() (any, error)
}

type file struct {
	path			string
	fileInfo		*fs.FileInfo
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewFile(path string) *file {
	r := file{
		path:		path,
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// FileIfc
// -------------------------------------------------------------------------------------------------

func (r *file) GetName() (*string, error) {
	fi, err := r.getFileInfo()
	if (nil != err) || (nil == fi) { return nil, err }
	v := (*fi).Name()
	return &v, nil
}

func (r *file) GetSize() (*int64, error) {
	fi, err := r.getFileInfo()
	if (nil != err) || (nil == fi) { return nil, err }
	v := (*fi).Size()
	return &v, nil
}

func (r *file) GetMode() (*fs.FileMode, error) {
	fi, err := r.getFileInfo()
	if (nil != err) || (nil == fi) { return nil, err }
	v := (*fi).Mode()
	return &v, nil
}

func (r *file) GetModTime() (*time.Time, error) {
	fi, err := r.getFileInfo()
	if (nil != err) || (nil == fi) { return nil, err }
	v := (*fi).ModTime()
	return &v, nil
}

func (r *file) IsDir() (bool, error) {
	fi, err := r.getFileInfo()
	if (nil != err) || (nil == fi) { return false, err }
	return (*fi).IsDir(), nil
}

func (r *file) GetSys() (any, error) {
	fi, err := r.getFileInfo()
	if (nil != err) || (nil == fi) { return nil, err }
	return (*fi).Sys(), nil
}

// -------------------------------------------------------------------------------------------------
// file
// -------------------------------------------------------------------------------------------------

// Get the FileInfo for this File with a read-through local cache copy
func (r *file) getFileInfo() (*fs.FileInfo, error) {
	// If we don't have a cached copy already
	if (nil == r.fileInfo) {
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

