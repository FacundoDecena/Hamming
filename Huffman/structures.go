package Huffman

import "github.com/pkg/errors"

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

//Example of use:

/*
	var item1, item2, item3 Huffman.Item
	item1.Symbol = 97
	item1.Weight = 2
	item2.Symbol = 98
	item2.Weight = 3
	item3.Symbol = 99
	item3.Weight = 4
	var treeNode1, treeNode2, treeNode3, tree *Huffman.TreeNode
	treeNode1, _ = treeNode1.New(item1)
	treeNode2, _ = treeNode2.New(item2)
	treeNode3, _ = treeNode3.New(item3)
	tree = tree.Insert(treeNode1, treeNode2)
	tree = tree.Insert(tree, treeNode3)
	fmt.Println(tree)
*/
//Resulting structure
/*
		                 (Symbol:0, Weight:9)
						/					\
(Symbol 99, Weight:4)						(Symbol: 0, Weight: 5)
/				\							/						\
nil				nil				(Symbol: 97, Weight: 2)				(Symbol: 98, Weight: 3)
								/					\				/					\
								nil					nil				nil					nil
*/
