package HuffmanCodification

type Aux struct {
	Caracter byte
	Length   int
}

func Deshuffman(bodyCoded []byte, table []byte) (originalBody []byte) {

	diccionary := make(map[uint32]Aux)
	var amountBytes uint32
	diccionary, amountBytes = stractTable(table)
	var integer uint32
	var result byte
	var numberOfShift int
	var length int
	var bitsTakenFromAByte int

	for index := 0; index < len(bodyCoded); index++ {
		bait := bodyCoded[index]
		for count := 0; count < 8; count++ {
			bitsTakenFromAByte++
			mask := doMask(bitsTakenFromAByte)
			baitAux := bait & mask
			baitAux = baitAux << uint(bitsTakenFromAByte-1)
			valor := uint(24 - numberOfShift)
			entero := uint32(baitAux)
			integer |= entero << valor // 24 or 16 or 8 or 0
			length++
			numberOfShift++
			if diccionary[integer].Length == length {
				// Element found.
				result = diccionary[integer].Caracter
				originalBody = append(originalBody, result)
				bait = bait << uint(bitsTakenFromAByte)
				bitsTakenFromAByte = 0
				integer = 0
				length = 0
				numberOfShift = 0

			}
		}
		bitsTakenFromAByte = 0

	}
	if len(originalBody) != int(amountBytes) {
		originalBody = originalBody[0 : len(originalBody)-1]

	}
	return originalBody

}
func stractTable(table []byte) (map[uint32]Aux, uint32) {
	var arrByte []byte
	var code uint32
	diccionary := make(map[uint32]Aux)
	var index int
	var amountBytes uint32
	arrByte = append(arrByte, table[0])
	arrByte = append(arrByte, table[1])
	arrByte = append(arrByte, table[2])
	arrByte = append(arrByte, table[3])

	amountBytes = uint32(arrByte[0])<<24 + uint32(arrByte[1])<<16 + uint32(arrByte[2])<<8 + uint32(arrByte[3])

	arrByte = nil

	for index = 4; index < len(table); index += 6 {
		var aux Aux
		aux.Caracter = table[index] // get the symbol

		// making de uint with the following 4 elements  after the first element of the array
		arrByte = append(arrByte, table[index+1])
		arrByte = append(arrByte, table[index+2])
		arrByte = append(arrByte, table[index+3])
		arrByte = append(arrByte, table[index+4])
		code = (uint32(arrByte[0]) << 24) + (uint32(arrByte[1]) << 16) + (uint32(arrByte[2]) << 8) + uint32(arrByte[3])
		// making the table.
		aux.Length = int(table[index+5])
		diccionary[code] = aux
		arrByte = nil
	}
	return diccionary, amountBytes
}
