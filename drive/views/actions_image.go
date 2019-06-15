package drivehandler

import (
	"drive/domain"
	"drive/drive"
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
	File *drive.File
	User *domain.Account
}

type ImageHandler struct {
	File     *drive.File
	User     *domain.Account
	Image    *drive.Image
	Siblings *drive.Siblings
}

func (v *ImageView) GetAction(r *http.Request, w http.ResponseWriter) error {
	enableCors(&w)
	v.ActionResponder.TemplateResponder.Template = "drive"
	file, user := v.File, v.User

	fmt.Printf("GetAction => Image %s\n", file.Name)

	if !file.Permissions.Read {
		return errors.New(errors.PermissionDenied, "Missing read permission for %v (User %s)", file.Path, user.Name)
	}

	image, err := drive.NewImage(v.File, v.User)
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
	image, err := drive.NewImage(v.File, v.User)
	if err != nil {
		return errors.Wrap(err, "NewImage")
	}

	image.Title = r.FormValue("title")
	image.Caption = r.FormValue("caption")
	image.Cutline = r.FormValue("cutline")
	if err := image.WriteMeta(v.User); err != nil {
		fmt.Println("image.Post", err)

		return errors.Wrap(err, "Image.WriteMeta")
	}
	fmt.Println("POST: no error", image)
	http.Redirect(w, r, v.File.Path, http.StatusFound)
	return nil
}

func (v *ImageView) DeleteAction(r *http.Request, w http.ResponseWriter) error {

	fmt.Printf("DeleteAction => File %s\n", v.File.Name)

	err := drive.DeleteFile(v.File)
	if err != nil {
		return errors.Wrap(err, "Could not delete file")
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
