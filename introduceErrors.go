package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func IntroduceErrors() {
	var fileName string
	var body []byte
	var err error
	r := bufio.NewReader(os.Stdin)

	clearScreen()
	fmt.Println("Ingrese el nombre del archivo a introducir errores CON EXTENCION")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Clean buffer
	_, _ = fmt.Fscanf(r, "%s")

	body, err = loadFile(fileName)
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	//Split the string between name and extention
	extension := strings.Split(fileName, ".")

	errorsFactory(body, extension[1])
}

func errorsFactory(file []byte, extension string) (fileWithErrors []byte) {
	switch extension {
	case "ha1":
		fileWithErrors = insertError(file, 7)
	case "ha2":
	case "ha3":
	case "ha4":
	}
	return file
}

func insertError(file []byte, kind int) []byte {
	return file
}
