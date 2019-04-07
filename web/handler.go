package web

import (
	"drive/session"
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
		rnd.HTML(w, http.StatusOK, "index", map[string]interface{}{"User": usr})
		return
	}
	if _, err := os.Stat(path.Join("./_static", r.URL.Path)); err != nil {
		//NotFound(w, r) os.IsNotExist(err)
		ErrorResponder2(w, err)
		return
	}
	http.FileServer(http.Dir("./_static")).ServeHTTP(w, r)
}

func assetHandler(prefix, location string) http.Handler {
	//fmt.Println("assetHandler", prefix, location)
	return http.StripPrefix(fmt.Sprintf("/%s/", prefix), http.FileServer(http.Dir(fmt.Sprintf("./%s", location))))
}
