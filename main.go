package main

import (
	"drive/arbeit"
	"drive/arbeit/pg_arbeit"
	"drive/config"
	"drive/drive/storage"
	drivehandler "drive/drive/views"
	"drive/web"
	"fmt"
	"path"
	"strings"
)

func main() {
	settings := config.ParseArgs()

	repo, err := pg_arbeit.GetDatabaseHandle()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer repo.Close()
	arbeit.Repo = repo

	//storage.SetDefaultStorage(config.Root)

	webserver := web.NewServer(config.Settings.Address.String())
	for name, stor := range settings.Storages {
		st := storage.Register(name, stor.Root, stor.BaseURL, stor.ServeURL, storage.GetGroupByID(stor.Group), stor.Mode)
		webserver.RegisterHandlerFunc(stor.ServeURL+"/", drivehandler.Serve(st))
		webserver.RegisterHandlerFunc(stor.BaseURL+"/", drivehandler.DispatchStorage(st))

	}
	webserver.RegisterHandlerFunc("/alben/", drivehandler.AlbumHandler(storage.Get("home"), "/alben"))
	//drivehandler.RegisterHandlers(webserver.RegisterHandlerFunc)
	arbeit.RegisterSubRouter(webserver.RegisterHandlerFunc)
	webserver.Run()
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

/*
package import
const var chan select
type struct interface map
defer
if else
switch case default fallthrough
for range break continue
goto
func return go


constants:
true false iota nil

types:
int int8 int16 int32 int64
uint uint8 uint16 uint32 uint64 uintptr
float32 float64 complex64 complex128
bool byte rune string error

functions:
make len cap new append copy close delete
complex real imag
panic recover
*/
