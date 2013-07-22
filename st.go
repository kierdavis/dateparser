package dateparser

import (
	"strings"
	"fmt"
)

type stnode struct {
	children map[rune]*stnode
	result int
}

type stinput struct {
	input string
	result int
}

func stBuild(inputs []stinput) (root *stnode) {
	root = &stnode{make(map[rune]*stnode), -1}
	
	for _, input := range inputs {
		node := root
		
		for _, c := range input.input {
			child, ok := node.children[c]
			if !ok {
				child = &stnode{make(map[rune]*stnode), -1}
				node.children[c] = child
			}
			
			node = child
		}
		
		node.result = input.result
	}
	
	return root
}

func (root *stnode) search(s string) (result int) {
	node := root
	
	for _, c := range strings.ToLower(s) {
		node = node.children[c]
		if node == nil {return -1}
	}
	
	return node.result
}

func (root *stnode) dump() {
	fmt.Printf("ROOT\n")
	root.dumpInternal("")
}

func (root *stnode) dumpInternal(prefix string) {
	i := 0
	
	for c, node := range root.children {
		if node.result == -1 {
			fmt.Printf("%s\\_ %c []\n", prefix, c)
		} else {
			fmt.Printf("%s\\_ %c [%d]\n", prefix, c, node.result)
		}
		
		if i == len(root.children) - 1 {
			node.dumpInternal(prefix + "   ")
		} else {
			node.dumpInternal(prefix + "|  ")
		}
		
		i++
	}
}
