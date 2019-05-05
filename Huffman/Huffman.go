package Huffman

//Item coso
type Item struct {
	Symbol byte
	Weight int
}

func toNode(table map[byte]int) (list []Item) {
	var newNode Item

	for key, value := range table {
		newNode.Symbol = key
		newNode.Weight = value

		list = append(list, newNode)
	}
	return list

}

func parva(list []Item) {

}
