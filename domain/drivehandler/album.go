package drivehandler

import (
	"drive/domain"
	"drive/domain/storage"
	"drive/domain/usecase"
	"drive/templates"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var rnd = templates.Rnd

type AlbumHandler struct {
	Folder *domain.Folder
	User   *domain.Account
	Images []*usecase.Image
}

func (h *AlbumHandler) Init(f *domain.File, u *domain.Account, s domain.Storage) {

}

func (a *AlbumHandler) Render(w http.ResponseWriter, r *http.Request) {

	path, _ := filepath.Rel("/alben", path.Clean(r.URL.Path))
	fmt.Printf(" - scanning '%s'\n", "/"+path)

	handle, err := storage.Get("public").GetHandle("/" + path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), 500)
	}
	file, _ := handle.ToFile(path, nil)
	if file.IsDir() {
		dir, _ := usecase.GetFolder(file, nil)
		_, _ = usecase.NewAlbum(dir)
		//album.Render(w, r)
		return
	} else {

		//		diary, _ := fs.NewDiary(file, storage)
		fmt.Println("DIARY")
		//diary.ServeHTTP(w, r)
	}

	contentType := r.Header.Get("Accept")

	if contentType == "application/json" {

		enc := json.NewEncoder(w)
		if err := enc.Encode(a); err != nil {
			log.Println(err)
		}
		// json, _ := json.Marshal(a)
		// w.Write(json)
	} else {
		fmt.Println("RENDER")
		err := rnd.HTML(w, http.StatusOK, "album", nil)
		if err != nil {
			fmt.Println("render error: ", err)
		}
	}
}
