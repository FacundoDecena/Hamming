package Huffman

import "container/heap"

func huffman(parva PriorityQueue) {

	var listMin []int
	heap.Init(&parva)

	// Print the order by Priority of expiry
	for parva.Len() > 0 {
		item := heap.Pop(&parva).(*Item)

		listMin = append(listMin, item.Weight)
		if len(listMin) == 2 {
			var newItem Item
			newItem = giveNewMin(listMin[0], listMin[1])
			// fmt.Printf("symbol: %c : Weight : %d \n", newItem.symbol, newItem.Weight)
			parva.Push(&newItem)

			listMin = nil
		}

		// use the tree with the item and make the tree.

	}

}

// This function take the map (table of frequencies) and make de list of items.

func toItems(table map[byte]int) (list []*Item) {

	for key, value := range table {
		var newNode Item
		newNode.Symbol = key
		newNode.Weight = value

		list = append(list, &newNode)
	}
	return list

}

// This function make the parva with a priority queue with the list of items.
func makeParva(listItems []*Item) PriorityQueue {

	priorityQueue := make(PriorityQueue, len(listItems))

	for i, item := range listItems {
		priorityQueue[i] = item
		priorityQueue[i].Index = i
	}

	return priorityQueue
}

func giveNewMin(min1 int, min2 int) (res Item) {
	res.Symbol = uint8(230)
	res.Weight = min1 + min2
	return res
}
