package Huffman

import "github.com/pkg/errors"

type TreeNode struct {
	left  *TreeNode
	value Item
	right *TreeNode
}

//Converts Item into a node for the huffman tree
func (treeNode *TreeNode) new(node Item) (*TreeNode, error) {
	if treeNode != nil {
		return nil, errors.New("treeNode must be null")
	}
	treeNode.left = nil
	treeNode.right = nil
	treeNode.value = node
	return treeNode, nil
}

//Inserts treeNodes to the tree
func (tree *TreeNode) insert(left TreeNode, right TreeNode) {
	var root Item
	root.symbol = -1
	root.weight = left.value.weight + right.value.weight

}
