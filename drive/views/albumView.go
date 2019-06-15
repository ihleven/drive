package drivehandler

import (
	"drive/drive"
	"drive/errors"
	"drive/session"
	"drive/templates"
	"fmt"
	"net/http"
	"path"
	"strings"
)

func AlbumListHandler(w http.ResponseWriter, r *http.Request) {

	//sessionUser, _ := session.GetSessionUser(r, w)

	err := templates.Render(w, http.StatusOK, "alben", map[string]interface{}{"Mallorca": "/Mallorca", "Hochyeitsreise": "/hochzeitsreise"})
	if err != nil {
		errors.Error(w, r, err)
	}
}

func AlbumHandler(storage drive.Storage, prefix string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		sessionUser, _ := session.GetSessionUser(r, w)

		folder, err := drive.GetFile(storage, strings.TrimPrefix(path.Clean(r.URL.Path), prefix), sessionUser)
		if err != nil {
			errors.Error(w, r, err)
			return
		}
		//if handle.Path() == "." {
		//	AlbumListHandler(w, r)
		//	return
		//}
		album, err := drive.NewAlbum(folder, sessionUser)
		if err != nil {
			errors.Error(w, r, err)
			return
		}

		data := map[string]interface{}{
			"album":   album,
			"folder":  folder,
			"account": sessionUser,
			"storage": storage,
		}
		fmt.Println("/////////", r.Header.Get("Accept"))
		switch r.Header.Get("Accept") {
		case "application/json":
			err = templates.SerializeJSON(w, http.StatusOK, data)
		default:
			err = templates.Render(w, http.StatusOK, "drive", data)
		}

		if err != nil {
			errors.Error(w, r, errors.Wrap(err, "render error"))
		}
		return

		if err != nil {
			errors.Error(w, r, err)
		}
	}
}
