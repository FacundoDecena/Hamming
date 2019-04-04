package main

import (
	"fmt"
	"io/ioutil"
	//"math/bits"
)

func main() {
	var fileName string
	maskFirst := 240
	maskLast := 15

	fmt.Println("Ingrese el nombre del archivo")
	_, _ = fmt.Scanf("%s", &fileName)

	body, err := loadFile(fileName)

	if err != nil {
		fmt.Println(err)
	} else {
		maskedFirst := body[0] & uint8(maskFirst)
		maskedLast := body[0] & uint8(maskLast)
		fmt.Printf("\nint: %08b masked first part: %08b masked last part: %08b", body[0], maskedFirst, maskedLast)
		xor := body[0] ^ body[1]
		fmt.Printf("\nbyte 1: %08b byte 2: %08b xor: %08b", body[0], body[1], xor)
	}
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

/*func saveFile(fileName string, body []byte) error{

	return ioutil.WriteFile(fileName, body, 0600)

}*/
