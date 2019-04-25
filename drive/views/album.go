package drivehandler

import (
	"drive/drive"
	"drive/drive/storage"
	"drive/session"
	"drive/templates"
	"drive/web"
	"net/http"
	"path"
	"path/filepath"
)

func AlbumHandler(w http.ResponseWriter, r *http.Request) {

	sessionUser, _ := session.GetSessionUser(r, w)

	path, _ := filepath.Rel("/alben", path.Clean(r.URL.Path))

	album, err := drive.GetAlbum(storage.Get("home"), "/"+path, sessionUser)
	if err != nil {
		web.Error(w, r, err)
		return
	}

	err = templates.Render(w, http.StatusOK, "album", map[string]interface{}{"Album": album})
	if err != nil {
		web.Error(w, r, err)
	}
}
