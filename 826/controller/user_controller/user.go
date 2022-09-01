package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/planet-i/go-example/826/modules/errors_modules"
	"github.com/planet-i/go-example/826/modules/user_modules"
	"github.com/planet-i/go-example/826/service/user_service"
)

func Get(c *gin.Context) error {
	var err error
	u := user_modules.User{}

	// check params
	err = c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "some message"})
		err = errors.Wrap(err, "request params bind error")
		return err
	}

	//params another check
	if u.ID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "id unavailable"})
		return errors_modules.ParamsErr
	}

	// check auth
	err = errors_modules.AuthError
	if err != nil {
		if err == errors_modules.AuthError {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "authorized failed"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "another failed"})
		}
		return err
	}

	// do something like check privilege
	err = errors_modules.PermissionError
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "no permission"})
		return err
	}

	user, err := user_service.Get(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "some error message"})
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
}
