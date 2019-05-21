package drivehandler

import (
	"drive/drive"
	"drive/errors"
	"drive/session"
	"drive/templates"
	"net/http"
	"path"
	"path/filepath"
)

func AlbumListHandler(w http.ResponseWriter, r *http.Request) {

	//sessionUser, _ := session.GetSessionUser(r, w)

	err := templates.Render(w, http.StatusOK, "alben", map[string]interface{}{"Mallorca": "/Mallorca", "Hochyeitsreise": "/hochzeitsreise"})
	if err != nil {
		errors.Error(w, r, err)
	}
}

func AlbumHandler(storage drive.Storage) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		//func AlbumHandler(w http.ResponseWriter, r *http.Request) {

		sessionUser, _ := session.GetSessionUser(r, w)

		path, _ := filepath.Rel("/alben", path.Clean(r.URL.Path))

		if path == "." {
			AlbumListHandler(w, r)
			return
		}
		album, err := drive.GetAlbum(storage, path, sessionUser)
		if err != nil {
			errors.Error(w, r, err)
			return
		}

		err = templates.Render(w, http.StatusOK, "album", map[string]interface{}{"Album": album})
		if err != nil {
			errors.Error(w, r, err)
		}
	}
}
