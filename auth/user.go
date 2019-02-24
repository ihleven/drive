package auth

import (
	"fmt"
	"os/user"
)

type Account struct {
	Uid, Gid      uint32
	Username      string
	Name          string
	HomeDir       string
	Authenticated bool
}

// An authentication backend is an interface with two required methods:
// get_user(user_id) and
// authenticate(request, **credentials),
// as well as a set of optional permission related authorization methods.

func Authenticate(username, password string) (*Account, error) {

	usr := &Account{0, 0, "anonymous", "Anonym", "", false}

	_, err := user.Lookup(username)
	if _, ok := err.(user.UnknownUserError); ok {
		fmt.Println("User not existing.")
		return nil, nil
	} else if err != nil {
		fmt.Println("Error looking up user.")
		return nil, err
	}

	//usr.Uid, usr.Gid = int(u.Uid), uint32(u.Gid)
	return &Account{501, 20, "mi", "Matthias Ihle", "", true}, nil

	switch username {
	case "mi":
		usr = &Account{501, 20, "mi", "Matthias Ihle", "", true}
	case "ihle":
		usr = &Account{1406, 1407, "ihle", "Matthias Ihle", "", true}
	default:

	}
	return usr, nil
}

var users = map[string]*user.User{}

func GetUserByID(uid string) *user.User {
	usr, ok := users[uid]
	if !ok {
		usr, err := user.LookupId(uid)
		if err != nil {
			return nil
		}
		users[uid] = usr
	}

	//switch err.(type) {
	//case user.UnknownUserIdError:
	//	owner = &user.User{Username: "unknown"}
	//default:
	//	owner = &user.User{Username: "unknown"}
	//}

	return usr
}

var groups = map[string]*user.Group{}

func GetGroupByID(gid string) *user.Group {

	grp, ok := groups[gid]
	if !ok {
		fmt.Println("GetGroupByID: looking up group", gid)
		grp, err := user.LookupGroupId(gid)
		if err != nil {
			fmt.Println("GetGroupByID: error looking up group", gid, err)
			return nil
		}
		groups[gid] = grp
	}
	return grp
}
