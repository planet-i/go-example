package main

import (
	"fmt"

	"github.com/google/wire"
)

type Leaf struct {
	Name string
}

func NewLeaf(name string) Leaf {
	return Leaf{Name: name}
}
func (l Leaf) LeafName() {
	fmt.Println("Leaf name", l.Name)
}

type Branch struct {
	L Leaf
}

func NewBranch(l Leaf) Branch {
	return Branch{L: l}
}

type Root struct {
	B Branch
}

func NewRoot(b Branch) Root {
	return Root{B: b}
}
func (r Root) GetLeafName() {
	r.B.L.LeafName()
}

// CommInitRoot 普通方法初始化一个组合深的结构
func CommInitRoot(name string) Root {
	// 需要很多的初始化步骤
	leaf := NewLeaf(name)
	branch := NewBranch(leaf)
	root := NewRoot(branch)
	return root
}

// InitRoot 用wire方式initRoot
func InitRoot(name string) *Root {
	wire.Build(NewLeaf, NewBranch, NewRoot)
	return nil
}
