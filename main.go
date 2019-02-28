package main

import (
	"drive/app"
	"drive/config"
	"drive/file/storage"
	"path"
	"strings"
)

func main() {
	config.ParseFlags()
	storage.SetDefaultStorage(config.Root)
	app.CreateServer()
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
