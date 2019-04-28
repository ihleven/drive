package web

import (
	"drive/errors"
	"drive/session"
	"drive/templates"
	"fmt"
	"net/http"
	"os"
	"path"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("IndexHandler", r.URL.Path)
	if r.URL.Path == "/" {
		/* Index page */
		usr, _ := session.GetSessionUser(r, w)
		templates.Render(w, http.StatusOK, "index", map[string]interface{}{"User": usr})
		return
	}
	if _, err := os.Stat(path.Join("./_static", r.URL.Path)); err != nil {
		//NotFound(w, r) os.IsNotExist(err)
		errors.Error(w, r, err)
		return
	}
	// TODO: nicht jedes mal einen FileServer registrieren
	http.FileServer(http.Dir("./_static")).ServeHTTP(w, r)
}

func assetHandler(prefix, location string) http.Handler {
	//fmt.Println("assetHandler", prefix, location)
	return http.StripPrefix(fmt.Sprintf("/%s/", prefix), http.FileServer(http.Dir(fmt.Sprintf("./%s", location))))
}
