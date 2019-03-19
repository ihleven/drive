package storage

import (
	"syscall"
	"time"
)

func statAtime(st *syscall.Stat_t) *time.Time {
	atime := time.Unix(st.Atim.Unix())
	return &atime
}
func statCtime(st *syscall.Stat_t) *time.Time {
	ctime := time.Unix(st.Atim.Unix())
	return &ctime
}

func (fh *FileHandle) StatTimes() (atime, mtime, ctime time.Time, err error) {

	mtime = fh.ModTime()
	stat := fh.Sys().(*syscall.Stat_t)
	atime = time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	return
}
