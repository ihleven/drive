package file

import (
	"fmt"
	"syscall"
)

func (fh *FileHandle) Write(buffer []byte) (n int, err error) {
	n, err = fh.File.Write(buffer)
	return
}

func (fh *FileHandle) ReadDir() ([]FileHandle, error) {

	entries, err := fh.File.Readdir(-1)
	if err != nil {
		return nil, err
	}
	//sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	var handles = make([]FileHandle, len(entries))
	for index, entry := range entries {
		if entry.Name()[0] == '.' {
			// ignore all files starting with '.'
			continue
		}
		handles[index] = FileHandle{FileInfo: entry}
		stat, _ := entry.Sys().(*syscall.Stat_t) // _ ist ok und kein error
		fmt.Println(stat)
	}
	return handles, nil

	// file := FileFromInfo(info)

	// r, w, _ := file.GetPermissions(usr.Uid, usr.Gid)
	//file.Permissions = struct{ Read, Write bool }{Read: r, Write: w}
}

func (f *FileHandle) GetPermissions(uid, gid uint32) (r, w, x bool) { // => handle

	//fmt.Println(uid, gid)
	Stat := f.Sys().(*syscall.Stat_t)

	if Stat == nil {
		return false, false, false
	}

	switch {
	case Stat.Uid == uid:
		ur, uw, ux := f.UserPermissions()
		r = r || ur
		w = w || uw
		x = x || ux
		//fmt.Println(r, w, x, ur, uw, ux)
		fallthrough
	case Stat.Gid == gid:
		gr, gw, gx := f.GroupPermissions()
		r = r || gr
		w = w || gw
		x = x || gx
		//fmt.Println(r, w, x, gr, gw, gx)
		fallthrough
	default:
		or, ow, ox := f.OthersPermissions()
		r = r || or
		w = w || ow
		x = x || ox
		//fmt.Println(r, w, x, or, ow, ox)
	}
	return r, w, x
}
func (f *FileHandle) UserPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[1] == 'r', rwx[2] == 'w', rwx[3] == 'x'
}
func (f *FileHandle) GroupPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[4] == 'r', rwx[5] == 'w', rwx[6] == 'x'
}
func (f *FileHandle) OthersPermissions() (r, w, x bool) {
	rwx := f.Mode().String()
	return rwx[7] == 'r', rwx[8] == 'w', rwx[9] == 'x'
}
