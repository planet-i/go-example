package user_dao

import "github.com/planet-i/go-example/826/modules/user_modules"

func Get(id int64) (user_modules.User, error) {
	user := user_modules.User{
		ID:   id,
		Name: "mitaka",
	}
	var err error
	return user, err
}
