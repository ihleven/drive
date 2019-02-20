package auth

import "fmt"

type User struct {
	Uid, Gid      uint32
	Username      string
	Authenticated bool
}

// An authentication backend is an interface with two required methods:
// get_user(user_id) and
// authenticate(request, **credentials),
// as well as a set of optional permission related authorization methods.

func Authenticate(username, password string) (*User, error) {

	var user *User

	switch username {
	case "mi":
		user = &User{501, 20, "mi", true}
	case "ihle":
		user = &User{1406, 1407, "ihle", true}
		fmt.Println("user:", user)
	default:
		user = &User{0, 0, "anonymous", false}
	}
	return user, nil
}
