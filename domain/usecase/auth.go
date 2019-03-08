package usecase

import (
	"drive/domain"
	"fmt"
	"os/user"
)

var cachedUsers = map[string]*user.User{}
var cachedGroups = map[string]*user.Group{}

func GetUserByID(uid string) *user.User {
	usr, ok := cachedUsers[uid]
	if !ok {
		usr, err := user.LookupId(uid)
		if err != nil {
			return nil
		}
		cachedUsers[uid] = usr
	}

	//switch err.(type) {
	//case user.UnknownUserIdError:
	//	owner = &user.User{Username: "unknown"}
	//default:
	//	owner = &user.User{Username: "unknown"}
	//}

	return usr
}

func GetGroupByID(gid string) *user.Group {

	grp, ok := cachedGroups[gid]
	if !ok {
		fmt.Println("GetGroupByID: looking up group", gid)
		grp, err := user.LookupGroupId(gid)
		if err != nil {
			fmt.Println("GetGroupByID: error looking up group", gid, err)
			return nil
		}
		cachedGroups[gid] = grp
	}
	return grp
}

func Authenticate(username, password string) (*domain.Account, error) {

	usr := &domain.Account{0, 0, "anonymous", "Anonym", "", false}

	_, err := user.Lookup(username)
	if _, ok := err.(user.UnknownUserError); ok {
		fmt.Println("User not existing.")
		return nil, nil
	} else if err != nil {
		fmt.Println("Error looking up user.")
		return nil, err
	}

	switch username {
	case "mi":
		usr = &domain.Account{501, 20, "mi", "Matthias Ihle", "", true}
	case "ihle":
		usr = &domain.Account{1406, 1407, "ihle", "Matthias Ihle", "", true}
	default:
		usr = &domain.Account{501, 20, "mi", "Matthias Ihle", "", true}
	}
	return usr, nil
}
