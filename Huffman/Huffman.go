package Huffman

import "container/heap"

// Function huffman receives a priority queue and do a binary tree to make the huffman codification.

func huffman(parva PriorityQueue) {

	var listMin []*TreeNode
	var tree *TreeNode
	// var code []int32

	heap.Init(&parva)

	// Print the order by Priority of expiry
	for parva.Len() > 0 {
		item := heap.Pop(&parva).(*TreeNode)
		listMin = append(listMin, item)
		if len(listMin) == 2 {
			tree = tree.Insert(listMin[0], listMin[1])
			parva.Push(tree)
			listMin = nil
		}

	}

}

// This function take the map (table of frequencies) and make de list of tree nodes.

func toItems(table map[byte]int) (list []*TreeNode) {

	for key, value := range table {
		var newItem Item
		var newTreeNode *TreeNode
		newItem.Symbol = key
		newItem.Weight = value
		newTreeNode, _ = newTreeNode.New(newItem)

		list = append(list, newTreeNode)
	}
	return list

}

// This function make the parva with a priority queue with the list of items.

func makeParva(listItems []*TreeNode) PriorityQueue {

	priorityQueue := make(PriorityQueue, len(listItems))

	for i, item := range listItems {
		priorityQueue[i] = item
		priorityQueue[i].Value.Index = i
	}

	return priorityQueue
}
