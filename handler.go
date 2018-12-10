package main

import (
	"compress/gzip"
	"compress/zlib"
	"container/list"
	"drive/storage"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

func pathRequestHandler(w http.ResponseWriter, req *http.Request) {

	url := path.Clean(req.URL.Path)
	filepath := path.Join((*root_folder), url)

	fileInfo, error := os.Stat(filepath)
	if error != nil {
		http.Error(w, "*PathError: stat() failure => "+error.Error(), 500)
		return
	}

	var context storage.Filer

	if fileInfo.IsDir() { // If it's a directory, open it !
		context, _ = storage.NewDirectory(fileInfo, filepath)
	} else {
		file, _ := storage.NewFile(fileInfo, url)
		if req.Method == "POST" {
			body := req.FormValue("body")
			f := &storage.File{Path: filepath, Body: []byte(body)}
			err := f.Save()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//http.Redirect(w, r, "/view/"+title, http.StatusFound)
		}
		//file.load()
		context = file
	}

	contentType := req.Header.Get("Content-type")
	if contentType == "application/json" {
		context.MarschalJSON(w, req)
	} else {
		context.RenderHTML(w, req)
	}

}

func handlePath(w http.ResponseWriter, req *http.Request) {

	filepath := path.Join((*root_folder), path.Clean(req.URL.Path))

	fmt.Printf("\"%s %s %s\" \"%s\" \"%s\"\n",
		req.Method,
		req.URL.String(),
		req.Proto,
		req.Referer(),
		req.UserAgent()) // TODO: Improve this crappy logging

	// Opening the file handle
	handle, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the file.", 404)
		return
	}
	defer handle.Close()

	// Checking if the opened handle is really a file
	statinfo, err := handle.Stat()
	if err != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		return
	}

	if statinfo.IsDir() { // If it's a directory, open it !
		handleDirectory(handle, w, req)
	} else {
		handleFile(handle, w, req)
	}
}

func handleDirectory(f *os.File, w http.ResponseWriter, req *http.Request) {

	fileInfos, _ := f.Readdir(-1) // returns all the FileInfo from the directory in a single slice

	// First, check if there is any index in this folder.
	for _, val := range fileInfos {
		if val.Name() == "index.html" {
			//handleFile(path.Join(f.Name(), "index.html"), w, req)
			return
		}
	}

	// Otherwise, generate folder content.
	children_dir_tmp := list.New()
	children_files_tmp := list.New()

	for _, val := range fileInfos {
		if val.Name()[0] == '.' {
			continue
		} // Remove hidden files from listing

		if val.IsDir() {
			children_dir_tmp.PushBack(val.Name())
		} else {
			children_files_tmp.PushBack(val.Name())
		}
	}

	// And transfer the content to the final array structure
	children_dir := copyToArray(children_dir_tmp)
	children_files := copyToArray(children_files_tmp)

	tpl, err := template.ParseFiles(dirlisting_tpl)
	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		log.Print(err)
		return
	}

	data := dirlisting{Name: req.URL.Path, ServerUA: "ihle",
		Children_dir: children_dir, Children_files: children_files}

	err = tpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func handleFile(f *os.File, w http.ResponseWriter, req *http.Request) {

	statinfo, err := f.Stat()
	if err != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		return
	}

	if (statinfo.Mode() &^ 07777) == os.ModeSocket { // If it's a socket, forbid it !
		http.Error(w, "403 Forbidden : you can't access this resource.", 403)
		return
	}

	// Manages If-Modified-Since and add Last-Modified (taken from Golang code)
	if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && statinfo.ModTime().Unix() <= t.Unix() {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Last-Modified", statinfo.ModTime().Format(http.TimeFormat))

	// Content-Type handling
	query, err := url.ParseQuery(req.URL.RawQuery)

	if err == nil && len(query["dl"]) > 0 { // The user explicitedly wanted to download the file (Dropbox style!)
		w.Header().Set("Content-Type", "application/octet-stream")
	} else {
		// Fetching file's mimetype and giving it to the browser
		if mimetype := mime.TypeByExtension(path.Ext(f.Name())); mimetype != "" {
			w.Header().Set("Content-Type", mimetype)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
	}

	// Manage gzip/zlib compression
	output_writer := w.(io.Writer)

	is_compressed_reply := false

	if !is_compressed_reply {
		// Add Content-Length
		w.Header().Set("Content-Length", strconv.FormatInt(statinfo.Size(), 10))
	}

	// Stream data out !
	buf := make([]byte, min(fs_maxbufsize, statinfo.Size()))
	n := 0
	for err == nil {
		n, err = f.Read(buf)
		output_writer.Write(buf[0:n])
	}

	// Closes current compressors
	switch output_writer.(type) {
	case *gzip.Writer:
		output_writer.(*gzip.Writer).Close()
	case *zlib.Writer:
		output_writer.(*zlib.Writer).Close()
	}

	f.Close()
}

func copyToArray(src *list.List) []string {
	dst := make([]string, src.Len())

	i := 0
	for e := src.Front(); e != nil; e = e.Next() {
		dst[i] = e.Value.(string)
		i = i + 1
	}

	return dst
}
