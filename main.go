package main

import (
	"fmt"
	"io/ioutil"
	//"strconv"
	//"math/bits"
)

func main() {
	var fileName string

	fmt.Println("Ingrese el nombre del archivo")
	fmt.Scanf("%s", &fileName)

	body, err := loadFile(fileName)

	if err != nil {
		fmt.Println(err)
	} else {
		xor := body[0] ^ body[1]
		fmt.Printf("%08b %08b %08b", body[0], body[1], xor)
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
