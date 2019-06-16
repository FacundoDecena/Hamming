package HuffmanCodification

import (
	"github.com/pkg/errors"
	"strings"
)

//TreeNode its a binary tree
type TreeNode struct {
	Left  *TreeNode
	Value Item
	Right *TreeNode
}

//Converts Item into a node for the huffman tree
func (treeNode *TreeNode) New(node Item) (*TreeNode, error) {
	if treeNode != nil {
		return nil, errors.New("treeNode must be nil")
	}
	treeNode = &TreeNode{
		Left:  nil,
		Right: nil,
		Value: node,
	}
	return treeNode, nil
}

//insert Inserts two treeNodes to the tree.
//
//Returns a new treeNode with symbol = 0, the sum of the weights.
//
//Points to the sons, on the left the lightest
func (tree *TreeNode) Insert(son1 *TreeNode, son2 *TreeNode) *TreeNode {
	var root *TreeNode
	var item Item

	item.Symbol = 0
	item.Weight = son1.Value.Weight + son2.Value.Weight
	root, _ = root.New(item)

	//This condition may be avoided if the heap is well managed.
	if son1.Value.Weight < son2.Value.Weight {
		root.Left = son1
		root.Right = son2
	} else {
		root.Left = son2
		root.Right = son1
	}

	return root
}

// This function go through the binary tree making the huffman codification (to right is 1, to left is a 0).
func (tree *TreeNode) GenerateCodification(codification string, codifications []string) []string {
	if tree.Right == nil && tree.Left == nil {
		buffer := strings.Builder{}
		buffer.WriteByte(tree.Value.Symbol)
		codification = codification + "@@" + buffer.String()
		codifications = append(codifications, codification)
		return codifications
	}
	codificationLeft := codification + "0"
	codifications = tree.Left.GenerateCodification(codificationLeft, codifications)
	codification += "1"
	codifications = tree.Right.GenerateCodification(codification, codifications)
	return codifications
}
