package arbeit

import (
	"drive/errors"
	"drive/templates"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func respond(w http.ResponseWriter, r *http.Request, template string, data map[string]interface{}) {

	switch r.Header.Get("Accept") {
	case "application/json":
		err := templates.SerializeJSON(w, http.StatusOK, data)
		if err != nil {
			errors.Error(w, r, errors.Wrap(err, "Could not serialize data"))
		}
	default:
		err := templates.Render(w, http.StatusOK, template, data)
		if err != nil {
			errors.Error(w, r, errors.Wrap(err, "render error"))
		}
	}

}
