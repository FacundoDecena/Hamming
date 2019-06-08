package Huffman

import (
	"bufio"
	"fmt"
	"os"
)

type Aux struct {
	Caracter byte
	Length   int
}

func callDeshuffman() {
	var ret []byte

	var fileName string
	r := bufio.NewReader(os.Stdin)

	clearScreen()
	fmt.Println("Ingrese el nombre del archivo sin extension")
	_, _ = fmt.Fscanf(r, "%s", &fileName)

	body, err := loadFile(fileName + ".huf")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("El texto comprimido es : %b \n", body)

		ret = deshuffman(body, fileName+".dic")

		fileName2 := fileName + ".dhu"
		err = saveFile(fileName2, ret)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func deshuffman(bodyCoded []byte, fileName string) (originalBody [](byte)) {

	diccionary := make(map[uint32]Aux)
	var amountBytes uint32
	diccionary, amountBytes = stractTable(fileName)
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
				if diccionary[integer].Caracter != 0 {
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
		}
		bitsTakenFromAByte = 0

	}
	if len(originalBody) != int(amountBytes) {
		originalBody = originalBody[0 : len(originalBody)-1]

	}
	return originalBody

}
func stractTable(fileName string) (map[uint32]Aux, uint32) {
	var body []byte
	var arrByte []byte
	var code uint32
	diccionary := make(map[uint32]Aux)
	var index int
	var amountBytes uint32

	body, err := loadFile(fileName)
	if err != nil {
		fmt.Println(err)
	} else {
		arrByte = append(arrByte, body[0])
		arrByte = append(arrByte, body[1])
		arrByte = append(arrByte, body[2])
		arrByte = append(arrByte, body[3])

		amountBytes = uint32(arrByte[0])<<24 + uint32(arrByte[1])<<16 + uint32(arrByte[2])<<8 + uint32(arrByte[3])

		arrByte = nil

		for index = 4; index < len(body); index += 6 {
			var aux Aux
			aux.Caracter = body[index] // get the symbol

			// making de uint with the following 4 elements  after the first element of the array
			arrByte = append(arrByte, body[index+1])
			arrByte = append(arrByte, body[index+2])
			arrByte = append(arrByte, body[index+3])
			arrByte = append(arrByte, body[index+4])
			code = (uint32(arrByte[0]) << 24) + (uint32(arrByte[1]) << 16) + (uint32(arrByte[2]) << 8) + uint32(arrByte[3])
			// making the table.
			aux.Length = int(body[index+5])
			fmt.Printf("codificacion del elemento de la tabla %c con longitud %d %32b\n", aux.Caracter, aux.Length, code)
			diccionary[code] = aux
			arrByte = nil
		}

	}
	return diccionary, amountBytes
}
