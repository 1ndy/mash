package main

import "fmt"

type TreeNode struct {
	Value    DocKey
	Children []TreeNode
}

func (t *TreeNode) addChild(n DocKey) {
	t.Children = append(t.Children, TreeNode{n, make([]TreeNode, 0)})
}

func (t *TreeNode) isValidPath(path []string) bool {
	if len(path) == 1 && t.Value.Key == path[0] {
		return true
	} else if t.Value.Key == path[0] {
		valid := false
		for _, child := range t.Children {
			valid = valid || child.isValidPath(path[1:])
		}
		return valid
	}
	return false
}

func (t *TreeNode) getPathStartNode(path []string) DocKey {
	root := *t
	for _, key := range path {
		for _, child := range root.Children {
			if child.Value.Key == key {
				root = child
				// fmt.Println(root.Value)
			}
		}
	}
	return root.Value
}

func (t *TreeNode) printTree() {
	//fmt.Printf("root: %s\nChildren: ", t.Value.Key)
	for _, child := range t.Children {
		fmt.Print(child.Value)
	}
	fmt.Println()
	for _, child := range t.Children {
		child.printTree()
	}
}

func placeDocKey(root *TreeNode, n DocKey, spaceInterval int) {
	root_level := root.Value.Spaces
	if root_level+spaceInterval == n.Spaces {
		root.addChild(n)
	} else if root_level+spaceInterval < n.Spaces {
		placeDocKey(&root.Children[len(root.Children)-1], n, spaceInterval)
	}
}

func numTreesInInput(keys []DocKey) int {
	numTrees := 0
	for _, key := range keys {
		if key.Spaces == 0 {
			numTrees++
		}
	}
	return numTrees
}

func splitKeyListIntoTrees(keys []DocKey) [][]DocKey {
	minIndent := findMinimumIndent(keys)
	var indices []int
	for i, key := range keys {
		if key.Spaces == minIndent {
			indices = append(indices, i)
		}
	}
	indices = indices[1:]
	var treeNodeList [][]DocKey
	prev := 0
	for _, v := range indices {
		treeNodeList = append(treeNodeList, keys[prev:v])
		prev = v
	}
	treeNodeList = append(treeNodeList, keys[prev:])
	return treeNodeList
}

func buildTree(keys []DocKey) TreeNode {
	// fmt.Printf("Building Tree '%s'\n", keys[0].Key)
	spaceInterval := findSpacingInterval(keys)
	root := TreeNode{keys[0], make([]TreeNode, 0)}
	keys = keys[1:]

	for _, key := range keys {
		placeDocKey(&root, key, spaceInterval)
	}
	return root
}
