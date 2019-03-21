package storage

import (
	"time"
)

func statAtime(st interface{}) *time.Time {
	atime := time.Unix(0, 0)
	return &atime
}
func statCtime(st interface{}) *time.Time {
	atime := time.Unix(0, 0)
	return &atime
}
