package usecase

import (
	"drive/domain"
	"fmt"
	"os/user"
)

var cachedUsers = map[uint32]*domain.User{}
var cachedGroups = map[uint32]*domain.Group{}

func GetUserByID(uid uint32) *domain.User {
	usr, ok := cachedUsers[uid]
	if !ok {
		usr, err := user.LookupId(fmt.Sprintf("%d", uid))
		if err != nil {
			return nil
		}
		cachedUsers[uid] = &domain.User{Uid: usr.Uid, Gid: usr.Gid, Username: usr.Username, Name: usr.Name, HomeDir: usr.HomeDir}
	}

	//switch err.(type) {
	//case user.UnknownUserIdError:
	//	owner = &user.User{Username: "unknown"}
	//default:
	//	owner = &user.User{Username: "unknown"}
	//}

	return usr
}

func GetGroupByID(gid uint32) *domain.Group {

	grp, ok := cachedGroups[gid]
	if !ok {
		grp, err := user.LookupGroupId(fmt.Sprintf("%d", gid))
		if err != nil {
			fmt.Println("GetGroupByID: error looking up group", gid, err)
			return nil
		}
		cachedGroups[gid] = &domain.Group{Gid: grp.Gid, Name: grp.Name}
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
