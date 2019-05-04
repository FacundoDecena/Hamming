package Huffman

type Item struct {
	symbol byte
	weight int
}

func toNode(table map[byte]int) (list []Item) {
	var newNode Item

	for key, value := range table {
		newNode.symbol = key
		newNode.weight = value

		list = append(list, newNode)
	}
	return list

}

func parva(list []Item) {

}
