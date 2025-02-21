package main

import (
	"fmt"
)

type Message string

func NewMessage() Message {
	return Message("hi here!")
}

type Greeter struct {
	Message Message
}

func NewGreeter(m Message) *Greeter {
	return &Greeter{Message: m}
}

func (g *Greeter) Greet() Message {
	return g.Message
}

type Event struct {
	Greeter *Greeter
}

func NewEvent(greeter *Greeter) *Event {
	return &Event{Greeter: greeter}
}
func (e *Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	// message := NewMessage()
	// greeter := NewGreeter(message)
	// event := NewEvent(greeter)
	// event.Start()
	event := initializeEvent()
	event.Start()
}
