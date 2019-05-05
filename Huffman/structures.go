package Huffman

import (
	"github.com/pkg/errors"
	"strings"
)

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

func (tree *TreeNode) GenerateCodification(codification string, codifications []string) []string {
	if tree.Right == nil && tree.Left == nil {
		buffer := strings.Builder{}
		buffer.WriteByte(tree.Value.Symbol)
		codification = codification + buffer.String()
		codifications = append(codifications, codification)
		return codifications
	}
	codificationLeft := codification + "0"
	codifications = tree.Left.GenerateCodification(codificationLeft, codifications)
	codification += "1"
	codifications = tree.Right.GenerateCodification(codification, codifications)
	return codifications
}

//Example of use:
/*
    var item1, item2, item3, item4, item5 Huffman.Item
	var code []string
	var temp string
	item1.Symbol = 97
	item1.Weight = 2
	item2.Symbol = 98
	item2.Weight = 3
	item3.Symbol = 99
	item3.Weight = 6
	item4.Symbol = 100
	item4.Weight = 8
	item5.Symbol = 101
	item5.Weight = 10
	var treeNode1, treeNode2, treeNode3, treeNode4, treeNode5, tree *Huffman.TreeNode
	treeNode1, _ = treeNode1.New(item1)
	treeNode2, _ = treeNode2.New(item2)
	treeNode3, _ = treeNode3.New(item3)
	treeNode4, _ = treeNode4.New(item4)
	treeNode5, _ = treeNode5.New(item5)
	tree = tree.Insert(treeNode1, treeNode2)
	tree = tree.Insert(tree, treeNode3)
	tree = tree.Insert(tree, treeNode4)
	tree = tree.Insert(tree, treeNode5)
	fmt.Println(tree.GenerateCodification(temp,code))

*/
