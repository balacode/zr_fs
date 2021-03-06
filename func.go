// -----------------------------------------------------------------------------
// ZR Library - File System Package                              zr-fs/[func.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package fs

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/balacode/zr"
)

// # File Functions
//   DirExists(path string) bool
//   FileExists(path string) bool
//   IsFileExt(filename string, fileExts []string) bool
//   IsTextFile(filename string) bool
//   ReadFileLines(filename string) []string
//   WriteFileLines(filename string, lines []string) error

// -----------------------------------------------------------------------------
// # File Functions

// DirExists returns true if the directory/folder given by 'path' exists.
func DirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
} //                                                                   DirExists

// FileExists returns true if the file given by 'path' exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
} //                                                                  FileExists

// FlatZip compresses the files specified in 'fileNames' into the ZIP archive
// named 'zipName'. Does not create subfolders, hence a 'flat' archive.
func FlatZip(zipName string, fileNames []string) error {
	var archive *zip.Writer
	{
		file, err := os.Create(zipName)
		if err != nil {
			return err
		}
		archive = zip.NewWriter(file)
		defer file.Close()
		defer archive.Close()
	}
	for _, fileName := range fileNames {
		var err error
		var file *os.File
		{
			file, err = os.Open(fileName)
			if err != nil {
				return err
			}
			defer file.Close()
		}
		var header *zip.FileHeader
		{
			var info os.FileInfo
			info, err = file.Stat()
			if err != nil {
				return err
			}
			header, err = zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = filepath.Base(fileName)
			header.Method = zip.Deflate
		}
		var wr io.Writer
		wr, err = archive.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(wr, file)
		if err != nil {
			return err
		}
	}
	return nil
} //                                                                     FlatZip

// IsFileExt returns true if the specified 'filename' has a
// file extension listed in 'fileExts'. The file extensions
// in the list should not include '.'. For example:
// []string{"go", "txt", "log"} is a valid list of file
// extensions, but []string{".go", "*.txt", ".log"} is not.
func IsFileExt(filename string, fileExts []string) bool {
	// TODO: only change to lower case on Windows
	filename = strings.ToLower(filename)
	for _, ext := range fileExts {
		if strings.HasSuffix(filename, strings.ToLower("."+ext)) {
			return true
		}
	}
	return false
} //                                                                   IsFileExt

// IsTextFile returns true if the given file name
// represents a text file type. For example "readme.txt"
// returns true, while "image.png" returns false.
func IsTextFile(filename string) bool {
	return IsFileExt(filename, TextFileExts)
} //                                                                  IsTextFile

// ReadFileLines reads the specified filename and returns
// all the lines it contains in a string array.
func ReadFileLines(filename string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {
		myError("Failed reading", filename, "due to:", err)
		return []string{} // erv
	}
	ret := strings.Split(string(data), "\n")
	return ret
} //                                                               ReadFileLines

// WriteFileLines writes lines to filename.
// This function is mainly used for saving text files.
func WriteFileLines(filename string, lines []string) error {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return myError(zr.EInvalidArg, "^filename")
	}
	data := []byte(strings.Join(lines, "\n"))
	//
	// terminate the last line with a newline
	if data[len(data)-1] != '\n' {
		//
		// handle Windows-type line breaks
		if bytes.Index(data, []byte("\r\n")) != -1 {
			data = append(data, '\r')
		}
		data = append(data, '\n')
	}
	// save the file
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return myError("Failed writing", filename, "due to:", err)
	}
	return nil
} //                                                              WriteFileLines

// end
