// DigiStratum GoLib - File handling
package fileio

import(
	"os"
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
)

// Write the contents of a string to a file
func WriteFileString(path string, content *string) error {
	c := []byte(*content)
	return ioutil.WriteFile(path, c, 0644)
}

// Write the contents of a []byte to a file
func WriteFileBytes(path string, content *[]byte) error {
	return ioutil.WriteFile(path, *content, 0644)
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
	tbuf, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, fmt.Errorf("Error reading '%s': %s", path, err.Error())
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

// Copy a file (src) to the destination (dst)
func CopyFile(src, dst string) error {
	// Source must be a file
	if ! IsFile(src) {
		return fmt.Errorf(
			"Files.CopyFile(): src (%s) is not a file", src,
		)
	}

	// Destination must either be a file (to be replaced) or a dir (to drop the file into)
	var destPath string
	if IsFile(dst) {
		destPath = dst
	} else if IsDir(dst) {
		// Keep the source filename, just send it to a new destination dir
		srcDir := path.Dir(src)
		srcFile := src[len(srcDir):]
		destPath = dst + "/" + srcFile
	} else {
		return fmt.Errorf(
			"Files.CopyFile(): dst (%s) is neither a file, nor a dir", dst,
		)
	}

	// Do the actual file copying bits
	in, err := os.Open(src)
	if err != nil { return err }
	defer in.Close()
	out, err := os.Create(destPath)
	if err != nil { return err }
	defer func() {
		cerr := out.Close()
		if err == nil { err = cerr }
	}()
	if _, err = io.Copy(out, in); err != nil { return err }
	err = out.Sync()
	return err
}

func GetDirFilesBySuffix(dir, suffix string) (*[]string, error) {
	files := make([]string, 0)
	if ! IsDir(dir) { return nil, fmt.Errorf("Not a directory: %s", dir) }
	if err := filepath.Walk(dir,
			func (file string, f os.FileInfo, err error) error {
				if nil != err { return err }				// Fail!
				if ! IsFile(file) { return nil }			// No Match
				if ! strings.HasSuffix(file, suffix) { return nil }	// No Match
				files = append(files, file)				// Match!
				return nil
			},
		); nil != err { return nil, err }
        return &files, nil
}

