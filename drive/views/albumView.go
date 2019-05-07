package drivehandler

import (
	"drive/drive"
	"drive/drive/storage"
	"drive/errors"
	"drive/session"
	"drive/templates"
	"net/http"
	"path"
	"path/filepath"
)

func AlbumListHandler(w http.ResponseWriter, r *http.Request) {

	sessionUser, _ := session.GetSessionUser(r, w)

	_, _ = filepath.Rel("/alben", path.Clean(r.URL.Path))

	album, err := drive.GetFile(storage.Get("home"), "/alben", sessionUser)
	if err != nil {
		errors.Error(w, r, err)
		return
	}

	err = templates.Render(w, http.StatusOK, "album", map[string]interface{}{"Album": album})
	if err != nil {
		errors.Error(w, r, err)
	}
}

func AlbumHandler(w http.ResponseWriter, r *http.Request) {

	sessionUser, _ := session.GetSessionUser(r, w)

	path, _ := filepath.Rel("/alben", path.Clean(r.URL.Path))

	album, err := drive.GetAlbum(storage.Get("home"), filepath.Join("/home", path), sessionUser)
	if err != nil {
		errors.Error(w, r, err)
		return
	}

	err = templates.Render(w, http.StatusOK, "album", map[string]interface{}{"Album": album})
	if err != nil {
		errors.Error(w, r, err)
	}
}
