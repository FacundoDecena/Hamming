package Huffman

func frequncies(list []byte) map[byte]int { // return the table of frequencies.

	table := make(map[byte]int)

	for index := 0; index < len(list); index++ {

		if table[list[index]] == 0 {
			table[list[index]] = 1
		} else {
			table[list[index]] += 1
		}

	}
	return table

}
