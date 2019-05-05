package Huffman

import "container/heap"

func toItems(table map[byte]int) (list []*Item) {

	for key, value := range table {
		var newNode Item
		newNode.symbol = key
		newNode.weight = value

		list = append(list, &newNode)
	}
	return list

}

func makeParva(listItems []*Item) PriorityQueue {

	priorityQueue := make(PriorityQueue, len(listItems))

	for i, item := range listItems {
		priorityQueue[i] = item
		priorityQueue[i].Index = i
	}

	heap.Init(&priorityQueue)

	var listMin []int

	// Print the order by Priority of expiry
	for priorityQueue.Len() > 0 {
		item := heap.Pop(&priorityQueue).(*Item)

		listMin = append(listMin, item.weight)
		if len(listMin) == 2 {
			var newItem Item
			newItem = giveNewMin(listMin[0], listMin[1])
			// fmt.Printf("symbol: %c : Weight : %d \n", newItem.symbol, newItem.weight)
			priorityQueue.Push(&newItem)

			listMin = nil
		}
	}
	return priorityQueue
}

func giveNewMin(min1 int, min2 int) (res Item) {
	res.symbol = uint8(230)
	res.weight = min1 + min2
	return res
}
