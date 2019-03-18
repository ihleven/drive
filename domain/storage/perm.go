package storage

import (
	"path/filepath"
	"syscall"
)

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

type Perm struct {
	uid, gid uint32
	stat     syscall.Stat_t

	r, w, x bool
}

func (st *FileSystemStorage) GetPerm(path string, uid, gid uint32) (*Perm, error) { // => storage

	fullpath := filepath.Join(st.Root, path)
	perm := &Perm{uid: uid, gid: gid}
	//stat = info.Sys().(*syscall.Stat_t)
	if err := syscall.Stat(fullpath, &perm.stat); err != nil {
		return nil, err
	}

	var mode = perm.stat.Mode
	var r, w, x bool

	switch {
	case perm.stat.Uid == uid:
		r, w, x = (mode&(OS_READ<<6)) != 0, (mode&(OS_WRITE<<6)) != 0, (mode&(OS_EX<<6)) != 0
		//fmt.Println(r, w, x)
		fallthrough
	case perm.stat.Gid == gid:
		r, w, x = r || ((mode&(OS_READ<<3)) != 0), w || ((mode&(OS_WRITE<<3)) != 0), x || ((mode&(OS_EX<<3)) != 0)
		//fmt.Println(r, w, x)
		fallthrough
	default:
		r, w, x = r || ((mode&(OS_READ<<0)) != 0), w || ((mode&(OS_WRITE<<0)) != 0), x || ((mode&(OS_EX<<0)) != 0)
		///fmt.Println(r, w, x)
	}

	//if (flag == os.O_RDONLY && !r) || (flag == os.O_WRONLY && !w) {
	//	return nil, os.ErrPermission
	//}
	//fmt.Printf("RWX:     %t %t %t \n", r, w, x)

	perm.r, perm.w, perm.x = r, w, x
	return perm, nil
}
