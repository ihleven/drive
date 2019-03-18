package usecase

import (
	"drive/domain"
	"errors"
	"path/filepath"
)

func GetReadHandle(storage domain.Storage, path string, uid, gid uint32) (domain.Handle, error) {

	handle, err := storage.GetHandle(path)
	if err != nil {
		return nil, err
	}
	if !handle.HasReadPermission(uid, gid) {
		return nil, errors.New("Permission denied")
	}
	return handle, nil
}

func GetFile(storage domain.Storage, path string, usr *domain.Account) (*domain.File, error) {

	handle, err := storage.GetHandle(path)
	if err != nil {
		return nil, err
	}

	file, err := handle.ToFile(path, usr)
	if err != nil {
		return nil, err
	}
	return file, nil
}
func GetFolder(storage domain.Storage, file *domain.File, usr *domain.Account) (*domain.Folder, error) {

	folder := &domain.Folder{File: file}
	//handles, err := file.ReadDirHandle()
	handles, err := storage.ReadDir(file.Path)
	if err != nil {
		return nil, err
	}

	for _, handle := range handles {

		entry, _ := handle.ToFile(filepath.Join(file.Path, handle.Name()), usr)

		_ = &domain.File{
			Handle: handle,
			Path:   filepath.Join(file.Path, handle.Name()),

			Name: handle.Name(),

			Mode:     handle.Mode(),
			Modified: handle.ModTime(),
			MIME:     handle.GuessMIME(),
			Owner:    &domain.User{},
			Group:    &domain.Group{},
		}
		folder.Entries = append(folder.Entries, entry)
		if entry.Name == "index.html" {
			folder.IndexFile = entry
		}
	}
	return folder, nil
}
