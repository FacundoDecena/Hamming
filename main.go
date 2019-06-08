package main

import (
	"./HammingCodification"
	"./Huffman"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var mainOp int
	r := bufio.NewReader(os.Stdin)
	continue_ := true
	for continue_ {
		clearScreen()
		fmt.Println("1 - Proteger archivo")
		fmt.Println("2 - Desporteger archivo")
		fmt.Println("3 - Introducir errores")
		fmt.Println("4 - Desproteger sin corregir errores")
		fmt.Println("5 - Ver detalles de archivos")
		fmt.Println("6 - Aplicar Huffman a un archivo")
		fmt.Println("7 - Salir")
		mainOp = 0
		_, _ = fmt.Fscanf(r, "%d", &mainOp)
		switch mainOp {
		case 1:
			Hamming()
		case 2:
			DeHamming(true)
		case 3:
			IntroduceErrors()
		case 4:
			DeHamming(false)
		case 5:
			seeSize()
		case 6:
			Huffman.Huffman()
		case 7:
			continue_ = false
		}
	}
}

func Hamming() {
	var dhOp int
	r := bufio.NewReader(os.Stdin)
	dhContinue_ := true
	for dhContinue_ {
		clearScreen()
		fmt.Println("¿Que tipo de PracticoDeMaquina quiere aplicar?")
		fmt.Println("1 - Hamming 7")
		fmt.Println("2 - Hamming 32")
		fmt.Println("3 - Hamming 1024")
		fmt.Println("4 - Hamming 32768")
		fmt.Println("5 - Volver")
		fmt.Printf("Su opcion: ")
		dhOp = 0
		_, _ = fmt.Fscanf(r, "%d", &dhOp)
		switch dhOp {
		case 1:
			preHamming(7)
		case 2:
			preHamming(32)
		case 3:
			preHamming(1024)
		case 4:
			preHamming(32768)
		case 5:
			dhContinue_ = false
		}
		_, _ = fmt.Fscanf(r, "%d")
	}

}

func preHamming(size int) {
	var fileName string
	var encodedBody []byte
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	fmt.Println("Ingrese el nombre del archivo con extensión.")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Since golang does not show the time a program runs...
	start := time.Now()
	body, err := loadFile(fileName)
	var fileType string
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		_, _ = fmt.Fscanf(r, "%s")
		return
	} else {
		switch size {
		case 7:
			fileType = ".ha1"
		case 32:
			fileType = ".ha2"
		case 1024:
			fileType = ".ha3"
		case 32768:
			fileType = ".ha4"
		}
		start = time.Now()
		if len(body) == 0 {
			encodedBody = []byte{}
		} else {
			switch size {
			case 7:
				encodedBody = HammingCodification.Hamming7(body)
			case 32:
				encodedBody = HammingCodification.Hamming(size, body)
			case 1024:
				encodedBody = HammingCodification.Hamming(size, body)
			case 32768:
				encodedBody = HammingCodification.Hamming(size, body)
			}
		}
		fileName = strings.Replace(fileName, ".txt", fileType, -1)
		err = saveFile(fileName, encodedBody)
		if err != nil {
			fmt.Println(err)
		}

	}
	elapsed := time.Since(start)
	log.Printf("\nHamming took %s", elapsed)
	_, _ = fmt.Fscanf(r, "%s")
	_, _ = fmt.Fscanf(r, "%s")
}

func DeHamming(fixErrors bool) {
	var dhOp int
	r := bufio.NewReader(os.Stdin)
	dhContinue_ := true
	for dhContinue_ {
		clearScreen()
		fmt.Println("¿Que tipo de PracticoDeMaquina ha sido aplicado?")
		fmt.Println("1 - Hamming 7")
		fmt.Println("2 - Hamming 32")
		fmt.Println("3 - Hamming 1024")
		fmt.Println("4 - Hamming 32768")
		fmt.Println("5 - Volver")
		fmt.Printf("Su opcion: ")
		dhOp = 0
		_, _ = fmt.Fscanf(r, "%d", &dhOp)
		switch dhOp {
		case 1:
			preDeHamming(7, fixErrors)
		case 2:
			preDeHamming(32, fixErrors)
		case 3:
			preDeHamming(1024, fixErrors)
		case 4:
			preDeHamming(32768, fixErrors)
		case 5:
			dhContinue_ = false
		}
		_, _ = fmt.Fscanf(r, "%d")
	}
}

func preDeHamming(size int, fixErrors bool) {
	var fileName string
	var body []byte
	var err error
	var start time.Time
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	var hammingCase string
	switch size {
	case 7:
		hammingCase = "1"
	case 32:
		hammingCase = "2"
	case 1024:
		hammingCase = "3"
	case 32768:
		hammingCase = "4"
	}
	fmt.Println("Ingrese el nombre del archivo .ha" + hammingCase + " o .he" + hammingCase + " con extension")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	extension := strings.Split(fileName, ".")
	if len(extension) >= 2 && (extension[1] == ("ha"+hammingCase) || extension[1] == ("he"+hammingCase)) {
		body, err = loadFile(fileName)
		if err != nil {
			fmt.Println(err)
			_, _ = fmt.Fscanf(r, "%d")
			_, _ = fmt.Fscanf(r, "%d")
			return
		}
		start = time.Now()
		var decodedFile []byte
		if len(body) == 0 {
			decodedFile = []byte{}
		} else {
			if size == 7 {
				decodedFile = HammingCodification.DeHamming7(body, fixErrors)
			} else {
				decodedFile = HammingCodification.CallDecode(size, body, fixErrors)
			}
		}
		if fixErrors {
			fileName = strings.Replace(fileName, "."+extension[1], ".dh"+hammingCase, -1)
		} else {
			fileName = strings.Replace(fileName, "."+extension[1], ".de"+hammingCase, -1)
		}
		err = saveFile(fileName, decodedFile)
		if err != nil {
			fmt.Println(err)
		}
		elapsed := time.Since(start)
		log.Printf("\nDeHamming took %s", elapsed)
		_, _ = fmt.Scanf("%s")
	} else {
		fmt.Println("La extension del archivo no es válida.")
		_, _ = fmt.Fscanf(r, "%s")
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
}

func IntroduceErrors() {
	var fileName string
	var body, fileWithErrors []byte
	var err error
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	fmt.Println("Ingrese el nombre del archivo a introducir errores extension:")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Clean buffer
	_, _ = fmt.Fscanf(r, "%s")
	body, err = loadFile(fileName)
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	//Split the string between name and extension
	extension := strings.Split(fileName, ".")
	switch extension[1] {
	case "ha1":
		fileWithErrors = HammingCodification.InsertError7(body)
	case "ha2":
		fileWithErrors = HammingCodification.InsertError(body, 32)
	case "ha3":
		fileWithErrors = HammingCodification.InsertError(body, 1024)
	case "ha4":
		fileWithErrors = HammingCodification.InsertError(body, 32768)
	default:
		fmt.Println("La extension del archivo no es válida.")
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	_ = saveFile(strings.Replace(fileName, ".ha", ".he", -1), fileWithErrors)
	fmt.Println("Se han introducido errores de manera correcta.")
	_, _ = fmt.Fscanf(r, "%s")
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Println("#####################################")
	fmt.Println("_______________HAMMING_______________")
	fmt.Println("#####################################")
	fmt.Println()
}
