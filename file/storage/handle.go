package storage

import (
	"errors"
	"fmt"
	"os"
)

type FileHandle struct {
	os.FileInfo
	File *os.File
	//Stat    *syscall.Stat_t
	storage Storage
}

func (fh *FileHandle) Close() {
	fh.File.Close()
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
