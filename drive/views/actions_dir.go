package drivehandler

import (
	"drive/domain"
	"drive/drive"
	"encoding/json"
	"fmt"
	"net/http"

	"drive/errors"
)

/*  **************
Dir
**************
*/

type DirActionResponder struct {
	ActionResponder
	File *drive.File
	User *domain.Account
}

func (a *DirActionResponder) GetAction(r *http.Request, w http.ResponseWriter) error {
	enableCors(&w)

	a.ActionResponder.TemplateResponder.Template = "drive"
	fmt.Printf("GetAction => Directory '%s'\n", a.File.Path)

	folder, err := drive.GetFolder(a.File, a.User)
	if err != nil {
		return errors.Wrap(err, "Error getFolder()")
	}

	a.Respond(w, r, map[string]interface{}{
		"Folder":      folder,
		"Account":     a.User,
		"Breadcrumbs": a.File.Breadcrumbs(),
	})

	return nil
}

func (a *DirActionResponder) PostAction(r *http.Request, w http.ResponseWriter) error {

	file, user := a.File, a.User
	fmt.Printf("PostAction => Directory \"%s/\"\n", file.Name)

	if !file.Permissions.Write {
		return errors.New(errors.PermissionDenied, "Missing write permissions for %s", file.Path)
	}

	handle, err := drive.UploadFile(file.Storage(), file.StoragePath(), r)
	if err != nil {
		return errors.Wrap(err, "Could not upload to folder '%v'", file.StoragePath())
	}
	//formfile, multipart, err := r.FormFile("file")
	//if err != nil {
	//	return errors.Wrap(err, "Failed to parse form ")
	//}
	//defer formfile.Close()

	//filename := filepath.Join(file.StoragePath(), multipart.Filename)
	//err = file.Storage().Save(filename, formfile, false)
	//if err != nil {
	//	return errors.Wrap(err, "Failed to save uploaded file %s", filename)
	//}

	a.Respond(w, r, map[string]interface{}{
		"File": handle.AsFile(),
		"User": user,
	})

	return nil
}

func (a *DirActionResponder) PutAction(r *http.Request, w http.ResponseWriter) error {

	file := a.File
	fmt.Printf("PutAction => Directory \"%s/\"\n", file.Name)

	//if !file.Permissions.Write {
	//	return errors.Errorf("no write permissions")
	//}

	var options struct {
		CreateThumbnails bool
	}
	err := json.NewDecoder(r.Body).Decode(&options)
	if err != nil {
		return errors.Wrap(err, "Error decoding put request body")
	}

	if options.CreateThumbnails {

		err := drive.MakeThumbs(file.Handle)
		if err != nil {
			return errors.Wrap(err, "Error making thumbnails")
		}
	}
	return nil
}
