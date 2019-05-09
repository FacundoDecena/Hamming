package Huffman

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//codificationAndLength
type CAL struct {
	Codification []byte
	Length       int
}

func Huffman() {
	var mainOp int
	r := bufio.NewReader(os.Stdin)
	continue_ := true
	for continue_ {
		clearScreen()
		fmt.Println("1 - Codificar")
		fmt.Println("2 - Decodificar")
		fmt.Println("3 - Salir")
		mainOp = 0
		_, _ = fmt.Fscanf(r, "%d", &mainOp)
		_, _ = fmt.Fscanf(r, "%s")
		switch mainOp {
		case 1:
			callHuffman()
			fmt.Println("Huffman aplicado correctamente")
			_, _ = fmt.Fscanf(r, "%s")
		case 2:
			callDeshuffman()
			fmt.Println("Huffman descomprimido correctamente")
			_, _ = fmt.Fscanf(r, "%s")
		case 3:
			continue_ = false
		}
	}
}

func callHuffman() {
	var fileName string

	r := bufio.NewReader(os.Stdin)

	clearScreen()
	fmt.Println("Ingrese el nombre del archivo")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Since golang does not show the time a program runs...

	body, err := loadFile(fileName)

	if err != nil {
		fmt.Println(err)
	} else {
		// Init Variables
		var listItems []*TreeNode
		var priorityQueue PriorityQueue
		var code []string

		table := make(map[byte]int)
		table = frequncies(body)
		listItems = toItems(table)
		priorityQueue = makeParva(listItems)

		code = huffman(priorityQueue)

		encodedBody, dictionary := encode(body, code)

		fileName = strings.Replace(fileName, "txt", "huf", -1)
		err = saveFile(fileName, encodedBody)
		if err != nil {
			fmt.Println(err)
			return
		}
		fileName = strings.Replace(fileName, "huf", "dic", -1)
		err = saveFile(fileName, dictionary)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
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

//encode applies Huffman
func encode(body []byte, code []string) (ret []byte, dic []byte) {
	//Create a dictionary
	var table map[byte]CAL
	var length int
	var tempCode byte
	//var codification []byte
	table = toMap(code)
	for k, v := range table {
		//Appends the key
		dic = append(dic, k)
		c := v.Codification
		for len(c) < 4 {
			c = append(c, 0)
		}
		//Append the values to the dictionary
		dic = append(dic, c...)
		dic = append(dic, byte(v.Length))
	}
	length = 0
	for i := 0; i < len(body); i++ {
		codificationI := table[body[i]].Codification
		lengthI := table[body[i]].Length
		//Compression
		for j := 0; j < len(codificationI); j++ {
			//Get the first byte
			codeJ := codificationI[j]
			//If the codification is not long enough keep going
			if length+lengthI < 8 {
				codeJ >>= uint(length)
				tempCode = tempCode | codeJ
				length += lengthI
			} else {
				//Save the part that fits in the byte
				firstPart := codeJ & ((exp(byte(length)) - 1) << uint(length))
				//Save the part that does not fit in the byte
				secondPart := codeJ & (exp(byte(length)) - 1)
				//Complete the byte
				tempCode = tempCode | (firstPart >> uint(length))
				//Save the completed byte to the ret structure
				ret = append(ret, tempCode)
				//Take the part that did not fit the byte
				tempCode = secondPart
				length = 8 - length - lengthI
			}
		}
	}
	if length != 0 {
		//Append the surplus
		ret = append(ret, tempCode)
	}

	return ret, dic
}

//toMap: from an easy to build structure to an easy to use structure
//
//gets a slice of strings. Each string consist of a symbol and its huffman codification.
//Returns a map with the symbols as keys and codifications with them length as values.
func toMap(table []string) map[byte]CAL {
	//ret := make(map[byte]uint32)
	ret := make(map[byte]CAL)

	for i := 0; i < len(table); i++ {
		var symbolString, codificationString string

		//Split the substrings
		fields := strings.Split(table[i], "@@")

		//strings.Fields separates the strings using white spaces, ignores the quantity
		//If the symbol is a white strings.Fields ignores it and we do not like that
		if len(fields) == 1 {
			symbolString = " "
			codificationString = fields[0]
		} else {
			//First string is the symbol
			symbolString = fields[1]

			//The rest is the symbol's codification
			codificationString = fields[0]
		}

		//Parse the strings to int
		codification64, _ := strconv.ParseInt(codificationString, 2, 32)

		//Cut the codification to 32 bits
		codification := uint32(codification64)
		//Get the length of the codification
		length := len(codificationString)
		//Move the surplus 0 to the right
		codification <<= uint32(32 - length)

		//Make a slice of bytes for the encode for each byte of body
		bs := make([]byte, 5)
		//Split the int32 into bytes
		binary.BigEndian.PutUint32(bs, codification)

		bs = takeBitsHuffman(length, bs, 0)

		//symbol has to be a byte
		symbol := []byte(symbolString)

		ret[symbol[0]] = CAL{Codification: bs, Length: length}
	}

	return ret
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Println("#####################################")
	fmt.Println("_______________HUFFMAN_______________")
	fmt.Println("#####################################")
	fmt.Println()
}

func loadFile(fileName string) ([]byte, error) {
	var err error
	var body []byte
	body, err = ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func saveFile(fileName string, body []byte) error {
	return ioutil.WriteFile(fileName, body, 0600)
}

//takeButsDeHamming
//
//bits is the amount of bits you need.
//input is the original byte slice. An extra byte is required
//initialPosition are the left shift to apply
func takeBitsHuffman(bits int, input []byte, initialPosition int) []byte {
	aux := input[len(input)-1]
	input[len(input)-1] = byte(0)
	bytesQuantity := int(math.Ceil(float64(bits+initialPosition) / float64(8)))
	ret := make([]byte, bytesQuantity)
	if initialPosition == 0 {
		for i := 0; i < bytesQuantity; i++ {
			ret[i] = input[i]
		}
		if bits%8 != 0 {
			ret[bytesQuantity-1] &= doMask(bits % 8)
		}
	} else {
		garbage := byte(0)
		for i := 0; i < bytesQuantity; i++ {
			ret[i] = garbage | ((doMask(8-initialPosition) & input[i]) >> byte(initialPosition))
			garbage = ((doMask(initialPosition) >> byte(8-initialPosition)) & input[i]) << byte(8-initialPosition)
		}
		mask := (bits%8 + initialPosition) % 8
		if mask == 0 {
			mask = 8
		}
		ret[bytesQuantity-1] &= doMask(mask)
	}
	input[len(input)-1] = aux
	return ret
}

// This function make a mask to take bits from a byte (left to right).
func doMask(bits int) uint8 {
	if bits > 8 {
		return uint8(0)
	} else if bits < 0 {
		return uint8(0)
	} else {
		val_mask := math.Pow(2, float64(bits)) - 1
		mask := uint8(val_mask) << uint(8-bits)
		return mask
	}

}

func exp(exponent byte) (ret byte) {
	ret = 1
	var i byte
	for i = 0; i < exponent; i++ {
		ret *= 2
	}
	return ret
}
