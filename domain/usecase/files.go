package usecase

import (
	"drive/domain"
)

func GetFile(storage, path string) (*domain.File, error) {

	file := storage.Open()
}

func GetFileForServing(storage, path string) (*domain.File, error) {

	file := Storage.Open()
}

type Storage interface {
	Open(name string) (*domain.File, error)
}
