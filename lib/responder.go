package lib

import (
	"drive/errors"
	"drive/templates"
	"net/http"
)

// Responder
type Responder interface {
	Respond(http.ResponseWriter, *http.Request, map[string]interface{}) error
	//Render(http.ResponseWriter, int, map[string]interface{}) error
}

// Default impl to embed in
type TemplateResponder struct {
	Template string
}

func (resp *TemplateResponder) Respond(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (err error) {

	switch r.Header.Get("Accept") {
	case "application/json":
		err = templates.SerializeJSON(w, http.StatusOK, data)
	default:
		err = templates.Render(w, http.StatusOK, resp.Template, data)
	}

	if err != nil {
		errors.Error(w, r, errors.Wrap(err, "render error"))
	}
	return
}

//func (resp *TemplateResponder) Render(w http.ResponseWriter, status int, data map[string]interface{}) error {
//
//	return Render(w, status, resp.template, data)
//}
