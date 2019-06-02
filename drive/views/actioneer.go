package drivehandler

import (
	"drive/lib"
	"drive/templates"
	"fmt"
	"net/http"

	"drive/errors"
)

type Actioneer interface {
	GetAction(*http.Request, http.ResponseWriter) error
	PostAction(*http.Request, http.ResponseWriter) error
	PutAction(*http.Request, http.ResponseWriter) error
	DeleteAction(*http.Request, http.ResponseWriter) error
	lib.Responder
	//Respond(http.ResponseWriter, *http.Request, map[string]interface{}) error
	//Render(http.ResponseWriter, int, string, map[string]interface{}) error
}

// Default imp
type ActionResponder struct {
	lib.TemplateResponder
	template string
}

func (a *ActionResponder) GetAction(r *http.Request, w http.ResponseWriter) error {
	fmt.Println("GetAction")
	return nil
}
func (a *ActionResponder) PostAction(r *http.Request, w http.ResponseWriter) error {
	fmt.Println("PostAction")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
	return nil
}
func (a *ActionResponder) DeleteAction(r *http.Request, w http.ResponseWriter) error {
	fmt.Println("DeleteAction")
	return errors.New(errors.NotImplemented, "Method Delete not implemented")
}
func (a *ActionResponder) PutAction(r *http.Request, w http.ResponseWriter) error {
	fmt.Println("PutAction")
	return errors.New(errors.NotImplemented, "Method Put not implemented")
}
func (a *ActionResponder) Respond_unused(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (err error) {

	switch r.Header.Get("Accept") {
	case "application/json":
		err = templates.SerializeJSON(w, http.StatusOK, data)
	default:
		//err = rnd.HTML(w, http.StatusOK, a.template, data)
		err = a.Render_unused(w, http.StatusOK, a.template, data)
	}

	if err != nil {
		errors.Error(w, r, errors.Wrap(err, "render error"))
	}
	return
}
func (a *ActionResponder) Render_unused(w http.ResponseWriter, status int, template string, data map[string]interface{}) error {
	fmt.Println("respond", data["File"])

	return templates.Render(w, status, template, data)
}
