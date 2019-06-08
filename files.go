package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

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

func seeSize() {
	extensions := []string{".txt", ".ha1", ".ha2", ".ha3", ".ha4", ".huf"}
	var fileName string
	r := bufio.NewReader(os.Stdin)

	fmt.Println("Ingrese el nombre del archivo SIN EXTENSION ")
	_, _ = fmt.Fscanf(r, "%s", &fileName)

	for index := 0; index < len(extensions); index++ {
		body, err := loadFile(fileName + extensions[index])
		if err != nil {
			fmt.Print("\n", err)
		} else {
			switch extensions[index] {
			case ".txt":
				fmt.Print("El archivo inicial tiene un tamaño de:", len(body), " Bytes ", " o ", (len(body))/1024, " KB")
			case ".ha1":
				fmt.Print("\n\n PracticoDeMaquinaTI2019 7 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha2":
				fmt.Print("\n\n PracticoDeMaquinaTI2019 32 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha3":
				fmt.Print("\n\n PracticoDeMaquina 1024 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha4":
				fmt.Print("\n\n PracticoDeMaquina 32"+
					"768 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".huf":
				fmt.Print("\n\n huffman tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			}

		}

	}

	fmt.Println("\n\n Presione enter para continuar")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	_, _ = fmt.Fscanf(r, "%s", &fileName)

}
