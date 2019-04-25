package web

import (
	"drive/errors"
	"drive/session"
	"drive/templates"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	//_ = fmt.Sprintf("ERROR: %+v", err)
	msg := fmt.Sprintf("%+v", err)
	rootCause := errors.RootCause(err)
	code := errors.GetCode(err)
	if code == errors.NoCode {
		_, status := toHTTPError(errors.RootCause(err))
		code = errors.ErrorCode(status)
		//msg = smsg
	}
	status := code.HTTPStatusCode()

	accept := r.Header.Get("Accept")

	switch {
	case accept == "application/json":
		err = templates.SerializeJSON(w, status, msg)
	case strings.Contains(accept, "text/html"):
		//session.Set(r, w, "debug", true)

		data := make(map[string]interface{})
		data["msg"] = msg
		data["subtitle"] = ""
		data["errno"] = status
		data["rootCause"] = rootCause
		debug, err := session.Get(r, "debug")
		if err == nil {
			if debugBool, ok := debug.(bool); ok {
				data["debug"] = debugBool
			}
		}

		templates.Render(w, status, "error", data)
	default:
		http.Error(w, msg, status)

	}

}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	typ := fmt.Sprintf("%T", err)
	if typ != "" {
		return typ, http.StatusBadRequest
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
