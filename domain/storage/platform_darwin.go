package storage

import (
	"syscall"
	"time"
)

func statMtime(st *syscall.Stat_t) time.Time {
	return time.Unix(st.Mtimespec.Unix())
}

func statAtime(st *syscall.Stat_t) *time.Time {
	atime := time.Unix(st.Atimespec.Unix())
	return &atime
}
func statCtime(st *syscall.Stat_t) *time.Time {
	ctime := time.Unix(st.Ctimespec.Unix())
	return &ctime
}

func (fh *FileHandle) StatTimes() (atime, mtime, ctime time.Time, err error) {

	mtime = fh.ModTime()
	stat := fh.Sys().(*syscall.Stat_t)
	atime = time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
	ctime = time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	return
}
