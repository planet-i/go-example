package main

type Leaf struct {
	Name string
}

type Branch struct {
	L Leaf
}

type Root struct {
	B Branch
}

func NewLeaf(name string) Leaf {
	return Leaf{Name: name}
}
func NewBranch(l Leaf) Branch {
	return Branch{L: l}
}
func NewRoot(b Branch) Root {
	return Root{B: b}
}
