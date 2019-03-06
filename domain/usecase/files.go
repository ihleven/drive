package usecase

import (
	"drive/domain"
	"drive/domain/storage"
)

func GetHandle(path string) (domain.Handle, error) {

	st := storage.Get("public")
	file, err := st.Open(path)
	return file, err
}

func GetFile(prefix, path string) (*domain.File, error) {

	st := storage.Get(prefix)
	handle, err := st.Open(path)
	file := &domain.File{
		Handle: handle,
		Path:   path,

		Name:  handle.Name(),
		Size:  handle.Size(),
		Mode:  handle.Mode(),
		MTime: handle.ModTime(),
		MIME:  handle.GuessMIME(),
	}
	return file, err
}

func GetServeContentHandle(prefix, path string, uid, gid uint32) (domain.Handle, error) {

	st := storage.Get(prefix)
	file, err := st.PermOpen(path, 0, uid, gid)
	return file, err
}
