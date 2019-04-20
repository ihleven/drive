package drivehandler

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Actioneer interface {
	GetAction(*http.Request, http.ResponseWriter) error
	PostAction(*http.Request, http.ResponseWriter) error
	DeleteAction(*http.Request, http.ResponseWriter) error
	//}
	//type Responder interface {
	Respond(http.ResponseWriter, *http.Request, map[string]interface{}) error
	Render(http.ResponseWriter, int, string, map[string]interface{}) error
}

// Default impl
type ActionResponder struct {
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
	return nil
}
func (a *ActionResponder) Respond(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (err error) {

	switch r.Header.Get("Accept") {
	case "application/json":
		err = rnd.JSON(w, http.StatusOK, data)
	default:
		//err = rnd.HTML(w, http.StatusOK, a.template, data)
		err = a.Render(w, http.StatusOK, a.template, data)
	}

	if err != nil {
		err = errors.Wrap(err, "render error")
	}
	return
}
func (a *ActionResponder) Render(w http.ResponseWriter, status int, template string, data map[string]interface{}) error {

	return rnd.HTML(w, status, template, data)
}
