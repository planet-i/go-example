package main

import (
	"github.com/google/wire"
)

func InitRoot(name string) Root {
	wire.Build(NewLeaf, NewBranch, NewRoot)
	return Root{}
}
