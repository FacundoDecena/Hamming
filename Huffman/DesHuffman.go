package Huffman

import (
	"bufio"
	"fmt"
	"os"
)

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

func deshuffman(bodyCoded []byte, fileName string) (originalBody []byte) {

	diccionary := make(map[uint32]byte)
	diccionary = stractTable(fileName)
	var arrByte []byte
	var result byte

	for index := 0; index < len(bodyCoded); index++ {
		arrByte = append(arrByte, bodyCoded[index])
		if len(arrByte) == 4 {
			integer := (uint32(arrByte[0]) << 24) + (uint32(arrByte[1]) << 16) + (uint32(arrByte[2]) << 8) + uint32(arrByte[3])
			result = diccionary[integer]
			originalBody = append(originalBody, result)
			arrByte = nil

		}
	}
	return originalBody

}

func stractTable(fileName string) map[uint32]byte {
	var body []byte
	var arrByte []byte
	var symbol byte
	var code uint32
	diccionary := make(map[uint32]byte)
	var index int

	body, err := loadFile(fileName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("El texto de la tabla es : %c \n", body)
		for index = 0; index < len(body); index += 5 {
			symbol = body[index] // get the symbol

			// making de uint with the following 4 elements  after the first element of the array
			arrByte = append(arrByte, body[index+1])
			arrByte = append(arrByte, body[index+2])
			arrByte = append(arrByte, body[index+3])
			arrByte = append(arrByte, body[index+4])
			code = (uint32(arrByte[0]) << 24) + (uint32(arrByte[1]) << 16) + (uint32(arrByte[2]) << 8) + uint32(arrByte[3])
			fmt.Printf("%32b\n", code)
			// making the table.
			diccionary[code] = symbol
			arrByte = nil
		}

	}
	return diccionary
}
