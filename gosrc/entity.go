package main

import "github.com/eminetto/clean-architecture-go/pkg/entity"

type Album struct {
	ID   entity.ID `json:"id" bson:"_id,omitempty"`
	Name string    `json:"name" bson:"name,omitempty"`
	Path string    `json:"path" bson:"path"`
}

type Repository interface {
	FindAll() ([]*Album, error)
	Get() (*Album, error)
}
