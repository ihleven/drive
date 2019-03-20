package web

import (
	"fmt"
	"net/http"
	"os"
)

func ErrorResponder(w http.ResponseWriter, msg string, errno int) {
	data := make(map[string]interface{})
	data["msg"] = msg
	data["errno"] = errno
	rnd.HTML(w, errno, "error", data)
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
func HttpLogOnError(w http.ResponseWriter, err error, message string) bool {
	if err == nil {
		return false
	}
	msg, code := toHTTPError(err)
	if message != "" {
		msg = fmt.Sprintf("%s => %s", msg, message)
	}
	http.Error(w, msg, code)
	return true
}

func HandleError(w http.ResponseWriter, err error, message string) bool {
	if err == nil {
		return false
	}
	msg, code := toHTTPError(err)
	if message != "" {
		msg = fmt.Sprintf("%s => %s", msg, message)
	}
	http.Error(w, msg, code)
	return true
}
