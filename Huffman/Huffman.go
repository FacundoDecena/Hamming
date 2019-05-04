package Huffman

import (
	"fmt"
)

type node struct {
	symbol byte
	weight int
}

func toNode(table map[byte]int) (list []node) {
	var newNode node

	for key, value := range table {
		newNode.symbol = key
		newNode.weight = value

		list = append(list, newNode)
	}
	return list

}

func parva(list []node) {

}
