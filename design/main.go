package main

import (
	"fmt"
	"time"

	"github.com/planet-i/go-example/design/builder"
)

func main() {
	CreateDBPool()
	CreateHouse()
}

func CreateHouse() {
	fmt.Println("--------CreateHouse")
	normalBuilder := builder.GetBuilder("normal")
	director := builder.NewDirector(normalBuilder)
	normalHouse := director.BuildHouse()

	fmt.Printf("Normal House Door Type: %s\n", normalHouse.DoorType)
	fmt.Printf("Normal House Window Type: %s\n", normalHouse.WindowType)
	fmt.Printf("Normal House Num Floor: %d\n", normalHouse.Floor)

	iglooBuilder := builder.GetBuilder("igloo")
	director.SetBuilder(iglooBuilder)
	iglooHouse := director.BuildHouse()

	fmt.Printf("\nIgloo House Door Type: %s\n", iglooHouse.DoorType)
	fmt.Printf("Igloo House Window Type: %s\n", iglooHouse.WindowType)
	fmt.Printf("Igloo House Num Floor: %d\n", iglooHouse.Floor)
}

func CreateDBPool() {
	fmt.Println("--------CreateDBPool")
	dbPool, err := builder.Builder().DSN("localhost:3306").MaxOpenConn(50).MaxConnLifeTime(10 * time.Second).Build()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dbPool)
}
