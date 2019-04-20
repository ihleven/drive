package usecase

import (
	"drive/domain"
	"drive/errors"
	"path/filepath"
)

func GetReadHandle(storage domain.Storage, path string, uid, gid uint32) (domain.Handle, error) {

	handle, err := storage.GetHandle(path)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get file handle")
	}
	if !handle.HasReadPermission(uid, gid) {
		return nil, errors.New(403, "uid: %v, gid %v has not read permission for %v", uid, gid, path)
	}
	return handle, nil
}

func GetFile(storage domain.Storage, path string, usr *domain.Account) (*domain.File, error) {

	handle, err := storage.GetHandle(path)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get file handle for %s", path)
	}

	file, err := handle.ToFile(path, usr)
	if err != nil {
		return nil, errors.Wrap(err, "Could not transform handle %v to File", file)
	}
	return file, nil
}

func DeleteFile(file *domain.File) error {

	if !file.Permissions.Write {
		return errors.New(errors.PermissionDenied, "Missing write permission for %s", file.Path)
	}
	err := file.Storage().Delete(file.Path)
	if err != nil {
		return errors.Wrap(err, "File '%s' could not be deleted!", file.Path)
	}
	return nil
}

func GetFolder(file *domain.File, usr *domain.Account) (*domain.Folder, error) {

	folder := &domain.Folder{File: file}
	//handles, err := file.ReadDirHandle()
	handles, err := file.ListDirHandles(false)
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
