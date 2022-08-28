package user_service

import (
	"github.com/pkg/errors"
	"github.com/planet-i/go-example/826/dao/user_dao"
	"github.com/planet-i/go-example/826/modules/user_modules"
)

var (
	UserIsForbidden = errors.New("user_service: user is forbidden")
)

func Get(id int64) (user_modules.User, error) {
	user, err := user_dao.Get(id)
	if err != nil {
		return user_modules.User{}, err
	}
	return user, nil
}
