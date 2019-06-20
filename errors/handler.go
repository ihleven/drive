package errors

import (
	"drive/session"
	"drive/templates"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	//fmt.Printf("ERROR: %+v\n", err)

	var txt string
	st, ok := err.(*stacktrace)
	if ok {
		txt = st.message
	}
	msg := fmt.Sprintf("%+v", err)
	rootCause := RootCause(err)
	code := GetCode(err)
	if code == NoCode {
		_, status := toHTTPError(RootCause(err))
		code = ErrorCode(status)
	}
	status := code.HTTPStatusCode()

	data := map[string]interface{}{
		"status":    status,
		"message":   txt,
		"msg":       msg[len(txt)+1:],
		"rootCause": rootCause,
	}
	//data["msg"] = msg
	//data["subtitle"] = ""
	//data["errno"] = status
	//data["rootCause"] = rootCause
	//fmt.Printf("ERROR: %+v\n", data)
	accept := r.Header.Get("Accept")

	switch {
	case accept == "application/json":
		err = templates.SerializeJSON(w, status, data)
	case strings.Contains(accept, "text/html"):
		//session.Set(r, w, "debug", true)

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
