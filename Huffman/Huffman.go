package Huffman

import (
	"container/heap"
	"encoding/binary"
	"strconv"
	"strings"
)

func callHuffman(body []byte) {

	// Init Variables
	var listItems []*TreeNode
	var priorityQueue PriorityQueue
	var code []string

	table := make(map[byte]int)
	table = frequncies(body)
	listItems = toItems(table)
	priorityQueue = makeParva(listItems)

	code = huffman(priorityQueue)

	encode(body, code)

}

// Function huffman receives a priority queue and do a binary tree to make the huffman codification.

func huffman(parva PriorityQueue) (codification []string) {

	var listMin []*TreeNode
	var tree *TreeNode
	var code []string
	var temp string

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
	codification = tree.GenerateCodification(temp, code)
	return codification

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

func encode(body []byte, code []string) (ret []byte) {
	//Create a dictionary
	var table map[byte]uint32
	table = toMap(code)

	for i := 0; i < len(body); i++ {
		//Make a slice of bytes for the encode for each byte of body
		bs := make([]byte, 4)
		//the value for body[i](A byte) in the dictionary
		binary.BigEndian.PutUint32(bs, table[body[i]])
		ret = append(ret, bs...)
	}
	return ret
}

//toMap: from an easy to build structure to an easy to use structure
//
//gets a slice of strings. Each string consist of a symbol and its huffman codification.
//Returns a map with the symbols as keys and codifications as values.
func toMap(table []string) map[byte]uint32 {
	ret := make(map[byte]uint32)

	for i := 0; i < len(table); i++ {
		var symbolString, codificationString string

		//Split the substrings
		fields := strings.Fields(table[i])

		//strings.Fields separates the strings using white spaces, ignores the quantity
		//If the symbol is a white strings.Fields ignores it and we do not like that
		if len(fields) == 1 {
			symbolString = " "
			codificationString = fields[0]
		} else {
			//First string is the symbol
			symbolString = fields[0]

			//The rest is the symbol's codification
			codificationString = fields[1]
		}

		//Parse the strings to int
		codification64, _ := strconv.ParseInt(codificationString, 2, 32)

		//Cut the codification to 32 bits
		codification := uint32(codification64)
		//Get the length of the codification
		length := len(fields)
		//Move the surplus 0 to the right
		codification <<= uint32(32 - length - 1)

		//The same for the symbol
		symbolInt, _ := strconv.ParseInt(symbolString, 2, 8)

		//But it has to be a byte
		symbol := byte(symbolInt)

		ret[symbol] = codification
	}

	return ret
}
