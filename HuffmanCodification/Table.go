package HuffmanCodification

func frequncies(list []byte) [256]int { // return the table of frequencies.

	// table := make(map[byte]int)
	var frequncies [256]int

	for index := 0; index < len(list); index++ {
		frequncies[list[index]] += 1
	}
	return frequncies

}
