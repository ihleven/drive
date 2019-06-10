package drivehandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"drive/domain"
	"drive/drive"
	"drive/errors"
)

type FileActionResponder2 struct {
	ActionResponder
	Handle drive.Handle
	User   *domain.Account
}

func (a *FileActionResponder2) GetAction(r *http.Request, w http.ResponseWriter) error {
	a.template = "file"
	handle, user := a.Handle, a.User

	fmt.Printf("GetAction => File %s\n", handle.Name())
	if !handle.HasReadPermission(user.Uid, user.Gid) {
		return errors.New(403, "uid: %v, gid %v has not read permission for %v", user.Uid, user.Gid, handle)
	}

	content, err := handle.GetUTF8Content()
	if err != nil {
		return errors.Wrap(err, "File init")
	}
	file, err := handle.ToFile(user)
	a.Respond(w, r, map[string]interface{}{
		"storage": handle.Storage(),
		"File":    file,
		"User":    user,
		"Content": content,
		"Title":   strings.TrimSuffix(handle.Name(), filepath.Ext(handle.Name())),
	})

	return nil
}

/*  **************
Files
**************
*/

type FileActionResponder struct {
	ActionResponder
	File *drive.File
	User *domain.Account
}

func (a *FileActionResponder) GetAction(r *http.Request, w http.ResponseWriter) error {

	a.ActionResponder.TemplateResponder.Template = "file"
	file, user := a.File, a.User
	fmt.Printf("GetAction => File %s\n", file.Name)

	if !file.Permissions.Read {
		return errors.New(403, "User '%v' has not read permission for %v", user.Username, file.Path)
	}

	content, err := file.GetUTF8Content()
	if err != nil {
		return errors.Wrap(err, "File init")
	}

	a.Respond(w, r, map[string]interface{}{
		"File":    file,
		"User":    user,
		"Content": content,
		"Title":   strings.TrimSuffix(file.Name, filepath.Ext(file.Name)),
	})

	return nil
}

func (a *FileActionResponder) PostAction(r *http.Request, w http.ResponseWriter) error {

	if !a.File.Permissions.Write {
		return fmt.Errorf("no write permissions")
	}

	formfile, multipart, err := r.FormFile("file")
	if err != nil {
		return errors.Wrap(err, "parsing form")
	}
	defer formfile.Close()
	_ = multipart.Header.Get("Content-Type")

	data, err := ioutil.ReadAll(formfile)
	if err != nil {
		return errors.Wrap(err, "read form file")
	}

	err = a.File.SetUTF8Content(data)
	if err != nil {
		return errors.Wrap(err, "writing utf8 file")
	}

	http.Redirect(w, r, a.File.Path, http.StatusFound)
	return nil
}

func (a *FileActionResponder) PutAction(r *http.Request, w http.ResponseWriter) error {

	body, _ := ioutil.ReadAll(r.Body)
	//i.update(body)
	json, _ := json.Marshal(body)
	w.Write(json)
	return nil
}

func (a *FileActionResponder) DeleteAction(r *http.Request, w http.ResponseWriter) error {
	file := a.File
	storage := file.Storage()
	fmt.Printf("DeleteAction => File %s\n", file.Name)

	if !file.Permissions.Write {
		return errors.Wrap(nil, "Failed to delete file")
	}
	_ = storage.Delete(file.Path)

	w.WriteHeader(http.StatusNoContent)

	return nil
}
