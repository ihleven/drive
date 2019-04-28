package domain

type Account struct {
	Uid, Gid      uint32
	Username      string
	Name          string
	HomeDir       string
	Authenticated bool
}
