package drivehandler

import (
	"drive/domain"
	"drive/domain/usecase"
	"fmt"
	"net/http"

	"drive/errors"
)

type DriveHandler interface {
	http.Handler
	//Init() error
	Get(r *http.Request) error
	Post(r *http.Request) error
	Respond(w http.ResponseWriter, format string) error
	//Render(w http.ResponseWriter)
}

type ImageView struct {
	ActionResponder
	File *domain.File
	User *domain.Account
}

type ImageHandler struct {
	File     *domain.File
	User     *domain.Account
	Image    *usecase.Image
	Siblings *domain.Siblings
}

func (v *ImageView) GetAction(r *http.Request, w http.ResponseWriter) error {
	v.template = "image"
	file, user := v.File, v.User

	fmt.Printf("GetAction => Image %s\n", file.Name)

	if !file.Permissions.Read {
		return errors.New(403, "User %v has not read permission for %v", user.Name, file.Path)
	}

	image, err := usecase.NewImage(v.File, v.User)
	if err != nil {
		return errors.Wrap(err, "NewImage")
	}
	siblings, err := v.File.Siblings()
	if err != nil {
		return errors.Wrap(err, "siblings error")
	}

	v.Respond(w, r, map[string]interface{}{
		"File":     file,
		"User":     user,
		"Image":    image,
		"Siblings": siblings,
	})

	return nil
}

func (v *ImageView) PostAction(r *http.Request, w http.ResponseWriter) error {

	if err := r.ParseForm(); err != nil {
		return errors.Wrap(err, "parse form")
	}
	image, err := usecase.NewImage(v.File, v.User)
	if err != nil {
		return errors.Wrap(err, "NewImage")
	}

	image.Title = r.FormValue("title")
	image.Caption = r.FormValue("caption")
	image.Cutline = r.FormValue("cutline")
	if err := image.WriteMeta(v.User); err != nil {
		return errors.Wrap(err, "Image.WriteMeta")
	}
	fmt.Println("POST: no error")
	http.Redirect(w, r, v.File.Path, http.StatusFound)
	return nil
}

func (v *ImageView) DeleteAction(r *http.Request, w http.ResponseWriter) error {

	fmt.Printf("DeleteAction => File %s\n", v.File.Name)

	err := usecase.DeleteFile(v.File)
	if err != nil {
		return errors.Wrap(err, "Could not delete file")
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
