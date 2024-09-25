package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	fmt.Println(MyUUID())
}

func MyUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
