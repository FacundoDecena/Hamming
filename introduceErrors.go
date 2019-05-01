package main

import (
	"bufio"
	"fmt"
	"os"
)

func IntroduceErrors() {
	var fileName string
	var body []byte
	var err error
	r := bufio.NewReader(os.Stdin)

	clearScreen()
	fmt.Println("Ingrese el nombre del archivo a introducir errores CON EXTENCION")
	_, _ = fmt.Fscanf(r, "%s", &fileName)

	body, err = loadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	errorsFactory(body)
}

func errorsFactory(file []byte) {

}
