//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func initializeEvent() *Event {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return nil
}
