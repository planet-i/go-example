package main

import (
	"fmt"
	"log"

	"github.com/planet-i/go-example/826/controller"
	"github.com/planet-i/go-example/826/modules/errors_modules"

	"github.com/pkg/errors"
)

func main() {
	newErr := errors.Wrap(errors_modules.AuthError, "this is a test")
	if errors.Cause(newErr) == errors_modules.AuthError {
		fmt.Println("new Err is errors_modules.AuthError")
	}
	e := controller.Route()
	log.Fatal(e.Run(":8080"))
}
