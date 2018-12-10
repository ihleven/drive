package views

import (
	"net/http"
)

type Renderer interface {
	render(w http.ResponseWriter, req *http.Request)
}

type View struct {
	template string
}
