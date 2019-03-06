package storage

import (
	"drive/domain"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type FileHandle struct {
	os.FileInfo
	File *os.File
	//Stat    *syscall.Stat_t
	storage domain.Storage
}

func (fh *FileHandle) GetFile() *os.File {
	return fh.File
}

func (fh *FileHandle) Close() error {
	if fh.File != nil {
		err := fh.File.Close()
		fh.File = nil
		return err
	}
	return errors.New("nil file")
}

func (fh *FileHandle) FormattedMTime() string {

	return fh.ModTime().Format(time.RFC822Z)
}

func (fh *FileHandle) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var buffer = make([]byte, fh.Size())
	fh.File.Seek(0, 0)
	len, err := fh.File.Read(buffer)
	if err != nil {
		return nil, err
	}

	if int64(len) != fh.Size() {
		return buffer, errors.New(fmt.Sprintf("read only %d of %d bytes", len, fh.Size()))
	}
	return buffer, nil
}

func (fh *FileHandle) SetContent(content []byte) error {

	return ioutil.WriteFile(fh.File.Name(), content, 0600)

	name := fh.File.Name()
	fd, err := os.OpenFile(name, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	n_bytes, err := fd.WriteString(strings.Replace(string(content), "\r\n", "\n", -1))
	if err != nil {
		fmt.Println("error writing file:", err, n_bytes, fh.File)
	}
	return err
}

func (fh *FileHandle) GetUTF8Content() (string, error) {

	content, err := fh.GetContent()
	if err != nil {
		return string(content), err
	}

	if utf8.Valid(content) {
		return string(content), nil
	} else {
		return string(content), errors.New("Invalid UTF-8")
	}
}
