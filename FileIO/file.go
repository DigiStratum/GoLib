package fileio

import (
	"os"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"time"
	"errors"
	"fmt"
)

type FileIfc interface {
	// Our own accessors
	GetPath() string			// Original path as supplied
	GetAbsPath() (*string, error)		// Absolute path, eliminates relative and more
	GetBasename() string			// Path without file name
	GetAbsBasename() (*string, error)	// Absolute path without file name

	Exists() bool				// Check if this file already exists on disk
	IsFile() (bool, error)			// Check if this is a regular file (vs. Dir, etc)

	CopyTo(path string) error		// Copy this file to another location

	ReadString() (*string, error)		// Read the file, return content as a *string
	ReadBytes() (*[]byte, error)		// Read the file, return contents as a *[]byte
	WriteString(content *string) error	// Write the contents of a string to a file
	WriteBytes(content *[]byte) error	// Write the contents of a []byte to a file

	// FileInfo accessors
	GetName() (*string, error)		// Name of the file without the path
	GetSize() (*int64, error)		// Size of the file as a count of bytes
	GetMode() (*fs.FileMode, error)		// File system "mode" (attributes)
	GetModTime() (*time.Time, error)	// Get the last modified timestamp
	IsDir() (bool, error)			// Check if this is a Dir (vs. regular File, etc)
	GetSys() (any, error)			// Get representation of data source (maybe nil!)
}

type file struct {
	path			string
	fileInfo		*fs.FileInfo
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewFile(path string) *file {
	// TODO: verify that path specifies a file, not a dir; reject dirs (use a different class for that!)
	r := file{
		path:		path,
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// FileIfc
// -------------------------------------------------------------------------------------------------

func (r *file) GetPath() string {
	return r.path
}

func (r *file) GetAbsPath() (*string, error) {
	absPath, err := filepath.Abs(r.path)
	if nil != err { return nil, err }
	return &absPath, nil
}

func (r *file) GetBasename() string {
	return filepath.Base(r.path)
}

func (r *file) GetAbsBasename() (*string, error) {
	// First get the absolute path
	absPath, err := filepath.Abs(r.path)
	if nil != err { return nil, err }
	// Then get the base of that
	basename := filepath.Base(absPath)
	return &basename, nil
}

func (r *file) Exists() bool {
	_, err := r.getFileInfo()
	return (nil == err) || ! errors.Is(err, os.ErrNotExist)
}

func (r *file) IsFile() (bool, error) {
	fi, err := r.getFileInfo()
	if (nil != err) || (nil == fi) { return false, err }
	v := (*fi).Mode()
	return v.IsRegular(), nil
}

func (r *file) CopyTo(path string) error {
	// Source must be a file
	if isFile, err := r.IsFile(); (! isFile) || (nil != err) {
		return fmt.Errorf("File.CopyTo(): src (%s) is not a file", r.path)
	}

	// Destination must either be a file (to be replaced) or a dir (to drop the file into)
	var destPath string
	destFile := NewFile(path)
	if ok, err := destFile.IsFile(); ok && (nil == err) {
		destPath = path
	} else if ok, err := destFile.IsDir(); ok && (nil == err) {
		// Keep the source filename, just send it to a new destination dir
		srcFile, err := r.GetName()
		if (nil != err) || (nil == srcFile) {
			return fmt.Errorf("File.CopyTo(): can't get filename from source path (%s)", r.path)
		}
		// TODO: only add PathSeparator if it's not already tacked onto path
		destPath = path + string(os.PathSeparator) + *srcFile
	} else {
		return fmt.Errorf(
			"File.CopyTo(): destination path (%s) is neither a file, nor a dir", destPath,
		)
	}

	// Do the actual file copying bits
	fin, err := os.Open(r.path)
	if err != nil { return err }
	defer fin.Close()
	fout, err := os.Create(destPath)
	if err != nil { return err }
	defer func() {
		cerr := fout.Close()
		if err == nil { err = cerr }
	}()
	if _, err = io.Copy(fout, fin); err != nil { return err }
	err = fout.Sync()
	return err
}

// Read the file located at the specified path and return the contents as a *string
func (r *file) ReadString() (*string, error) {
	tbuf, err := r.ReadBytes()
	if nil != err { return nil, err }
	s := string(*tbuf)
	return &s, nil
}

// Read the file located at the specified path and return the contents as a *[]byte
func (r *file) ReadBytes() (*[]byte, error) {
	tbuf, err := ioutil.ReadFile(r.path)
	if nil != err {
		return nil, fmt.Errorf("Error reading '%s': %s", r.path, err.Error())
	}
	return &tbuf, nil
}

func (r *file) WriteString(content *string) error {
        c := []byte(*content)
	err := ioutil.WriteFile(r.path, c, 0644)
	if nil == err { r.fileInfo = nil }
	return err
}

func (r *file) WriteBytes(content *[]byte) error {
	err := ioutil.WriteFile(r.path, *content, 0644)
	if nil == err { r.fileInfo = nil }
	return err
}

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

