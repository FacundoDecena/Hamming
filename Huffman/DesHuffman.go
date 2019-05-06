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
		fmt.Printf("El texto comprimido es : %c \n", body)

		ret = deshuffman(body, fileName+".dic")

		fmt.Printf("el texto es: %c\n\n", ret)
	}
}

func deshuffman(bodyCoded []byte, fileName string) (originalBody [](byte)) {

	diccionary := make(map[uint32]Aux)
	diccionary = stractTable(fileName)
	var integer uint32
	var result byte
	var numberOfShift int
	var length int

	for index := 0; index < len(bodyCoded); index++ {
		bait := bodyCoded[index]
		numberOfShift++
		for count := 0; count < 8; count++ {
			mask := doMask(count + 1)
			baitAux := bait & mask
			valor := uint(8 * (4 - (numberOfShift)))
			entero := uint32(baitAux)
			integer |= entero << valor // 24 or 16 or 8 or 0
			length++
			if diccionary[integer].Length == length {
				if diccionary[integer].Caracter != 0 {
					// Element found.
					result = diccionary[integer].Caracter
					originalBody = append(originalBody, result)
					integer = 0
					length = 0
					numberOfShift = 0
				}
			}
		}

	}
	return originalBody

}
func stractTable(fileName string) map[uint32]Aux {
	var body []byte
	var arrByte []byte
	var code uint32
	diccionary := make(map[uint32]Aux)
	var index int

	body, err := loadFile(fileName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("El texto de la tabla es : %c \n", body)
		for index = 0; index < len(body); index += 6 {
			var aux Aux
			aux.Caracter = body[index] // get the symbol

			// making de uint with the following 4 elements  after the first element of the array
			arrByte = append(arrByte, body[index+1])
			arrByte = append(arrByte, body[index+2])
			arrByte = append(arrByte, body[index+3])
			arrByte = append(arrByte, body[index+4])
			code = (uint32(arrByte[0]) << 24) + (uint32(arrByte[1]) << 16) + (uint32(arrByte[2]) << 8) + uint32(arrByte[3])
			fmt.Printf("codificacion del elemento de la tabla %32b\n", code)
			// making the table.
			aux.Length = int(body[index+5])
			diccionary[code] = aux
			arrByte = nil
		}

	}
	return diccionary
}
