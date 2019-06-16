package main

import (
	"./HammingCodification"
	"./HuffmanCodification"
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	var mainOp int
	r := bufio.NewReader(os.Stdin)
	continue_ := true
	for continue_ {
		clearScreen()
		fmt.Println("1 - Hamming")
		fmt.Println("2 - Huffman")
		fmt.Println("3 - Ver detalles de archivos")
		fmt.Println("4 - Salir")
		fmt.Print("Su opcion: ")
		mainOp = 0
		_, _ = fmt.Fscanf(r, "%d", &mainOp)
		switch mainOp {
		case 1:
			menuHamming()
		case 2:
			menuHuffman()
		case 3:
			statistics()
		case 4:
			continue_ = false
		}

	}
}

func menuHamming() {
	var mainOp int
	r := bufio.NewReader(os.Stdin)
	continue_ := true
	for continue_ {
		clearScreen()
		fmt.Println("1 - Proteger archivo")
		fmt.Println("2 - Desporteger archivo")
		fmt.Println("3 - Introducir errores")
		fmt.Println("4 - Desproteger sin corregir errores")
		fmt.Println("5 - Volver")
		fmt.Print("Su opcion: ")
		mainOp = 0
		_, _ = fmt.Fscanf(r, "%d", &mainOp)
		switch mainOp {
		case 1:
			hamming()
		case 2:
			deHamming(true)
		case 3:
			introduceErrors()
		case 4:
			deHamming(false)
		case 5:
			continue_ = false
		}
	}
}

func menuHuffman() {
	var mainOp int
	r := bufio.NewReader(os.Stdin)
	continue_ := true
	for continue_ {
		clearScreen()
		fmt.Println("1 - Codificar")
		fmt.Println("2 - Decodificar")
		fmt.Println("3 - Volver")
		fmt.Print("Su opcion: ")
		mainOp = 0
		_, _ = fmt.Fscanf(r, "%d", &mainOp)
		_, _ = fmt.Fscanf(r, "%s")
		switch mainOp {
		case 1:
			huffman()
		case 2:
			desHuffman()
		case 3:
			continue_ = false
		}
	}
}

func hamming() {
	var dhOp int
	r := bufio.NewReader(os.Stdin)
	dhContinue_ := true
	for dhContinue_ {
		clearScreen()
		fmt.Println("¿Que tipo de Hamming quiere aplicar?")
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
	}
}

func deHamming(fixErrors bool) {
	var dhOp int
	r := bufio.NewReader(os.Stdin)
	dhContinue_ := true
	for dhContinue_ {
		clearScreen()
		fmt.Println("¿Que tipo de Hamming ha sido aplicado?")
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
	}
}

func introduceErrors() {
	var fileName string
	var body, fileWithErrors []byte
	var err error
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	fmt.Println("Ingrese el nombre del archivo a introducir errores con extension:")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Clean buffer
	_, _ = fmt.Fscanf(r, "%s")
	originalText, err := loadFile(fileName, false)
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	//Split the string between name and extension
	extension := strings.Split(fileName, ".")
	switch extension[1] {
	case "ha1":
		body = originalText[:len(originalText)-10]
		fileWithErrors = append(HammingCodification.InsertError7(body), originalText[len(originalText)-10:]...)
	case "ha2":
		body = originalText[:len(originalText)-20]
		fileWithErrors = append(HammingCodification.InsertError(body, 32), originalText[len(originalText)-20:]...)
	case "ha3":
		body = originalText[:len(originalText)-20]
		fileWithErrors = append(HammingCodification.InsertError(body, 1024), originalText[len(originalText)-20:]...)
	case "ha4":
		body = originalText[:len(originalText)-20]
		fileWithErrors = append(HammingCodification.InsertError(body, 32768), originalText[len(originalText)-20:]...)
	default:
		fmt.Println("La extension del archivo no es válida.")
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	_ = saveFile(strings.Replace(fileName, ".ha", ".he", -1), fileWithErrors)
	fmt.Println("Se han introducido errores de manera correcta.")
	_, _ = fmt.Fscanf(r, "%s")
}

func preHamming(size int) {
	var fileName string
	var encodedBody []byte
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	unixDate, err := askDate()
	if err != nil {
		fmt.Print(err)
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	fmt.Println("Ingrese el nombre del archivo con extensión.")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Since golang does not show the time a program runs...
	start := time.Now()
	body, err := loadFile(fileName, false)
	var fileType string
	if err != nil {
		fmt.Println(err)
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
			if size != 7 {
				for i := 0; i < 10; i++ {
					encodedBody = append(encodedBody, byte(48))
				}
			}
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
		encodedBody = append(encodedBody, unixDate...)
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
		body, err = loadFile(fileName, true)
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

func statistics() {
	clearScreen()
	extensions := []string{".txt", ".ha1", ".ha2", ".ha3", ".ha4", ".huf"}
	var fileName string
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el nombre del archivo sin extension.")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	clearScreen()
	for index := 0; index < len(extensions); index++ {
		body, err := loadFile(fileName+extensions[index], false)
		if err != nil {
			fmt.Print("\n\n ", err)
		} else {
			switch extensions[index] {
			case ".txt":
				fmt.Print(" El archivo original tiene un tamaño de:", len(body), " Bytes ", " o ", (len(body))/1024, " KB")
			case ".ha1":
				fmt.Print("\n\n Hamming 7 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha2":
				fmt.Print("\n\n Hamming 32 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha3":
				fmt.Print("\n\n Hamming 1024 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha4":
				fmt.Print("\n\n Hamming 32"+
					"768 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".huf":
				fmt.Print("\n\n Huffman tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			}
		}
	}
	fmt.Println("\n\n Presione enter para continuar")
	_, _ = fmt.Fscanf(r, "%s")
	_, _ = fmt.Fscanf(r, "%s")
}

func huffman() {
	var fileName string
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	unixDate, err := askDate()
	if err != nil {
		fmt.Print(err)
		_, _ = fmt.Fscanf(r, "%s")
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	fmt.Println("Ingrese el nombre del archivo con extension.")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Since golang does not show the time a program runs...
	start := time.Now()
	body, err := loadFile(fileName, false)
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		_, _ = fmt.Fscanf(r, "%s")
		return
	} else {
		encodedBody, dictionary := HuffmanCodification.CallHuffman(body)
		dictionary = append(dictionary, unixDate...)
		fileName = strings.Split(fileName, ".")[0]
		fileName = fileName + ".huf"
		err = saveFile(fileName, encodedBody)
		if err != nil {
			fmt.Println(err)
			_, _ = fmt.Fscanf(r, "%s")
			_, _ = fmt.Fscanf(r, "%s")
			return
		}
		fileName = strings.Replace(fileName, "huf", "dic", -1)
		err = saveFile(fileName, dictionary)
		if err != nil {
			fmt.Println(err)
			_, _ = fmt.Fscanf(r, "%s")
			_, _ = fmt.Fscanf(r, "%s")
			return
		}
	}
	elapsed := time.Since(start)
	log.Printf("\nHuffman took %s", elapsed)
	_, _ = fmt.Fscanf(r, "%s")
	_, _ = fmt.Fscanf(r, "%s")
}

func desHuffman() {
	var fileName string
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	fmt.Println("Ingrese el nombre del archivo sin extension")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Since golang does not show the time a program runs...
	start := time.Now()
	body, err := loadFile(fileName+".huf", false)
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	table, err := loadFile(fileName+".dic", true)
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		_, _ = fmt.Fscanf(r, "%s")
		return
	} else {
		decodedBody := HuffmanCodification.Deshuffman(body, table)
		fileName = fileName + ".dhu"
		err = saveFile(fileName, decodedBody)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	elapsed := time.Since(start)
	log.Printf("\nDeshuffman took %s", elapsed)
	_, _ = fmt.Fscanf(r, "%s")
	_, _ = fmt.Fscanf(r, "%s")
}

func askDate() ([]byte, error) {
	//Ask for the date
	clearScreen()
	r := bufio.NewReader(os.Stdin)
	var auxDay, auxMonth, auxYear, auxHour, auxMinutes, auxSeconds string
	var day, month, year, hour, minutes, seconds int
	err := make([]error, 6)
	fmt.Println("Ingrese el dia, mes, año, hora, minutos y segundos en los que quiere la decodificacin del archivo este disponible:")
	fmt.Print("Dia: ")
	_, _ = fmt.Fscanf(r, "%s", &auxDay)
	_, _ = fmt.Fscanf(r, "%d")
	fmt.Print("Mes: ")
	_, _ = fmt.Fscanf(r, "%s", &auxMonth)
	_, _ = fmt.Fscanf(r, "%d")
	fmt.Print("Año: ")
	_, _ = fmt.Fscanf(r, "%s", &auxYear)
	_, _ = fmt.Fscanf(r, "%d")
	fmt.Print("Hora: ")
	_, _ = fmt.Fscanf(r, "%s", &auxHour)
	_, _ = fmt.Fscanf(r, "%d")
	fmt.Print("Minutos: ")
	_, _ = fmt.Fscanf(r, "%s", &auxMinutes)
	_, _ = fmt.Fscanf(r, "%d")
	fmt.Print("Segundos: ")
	_, _ = fmt.Fscanf(r, "%s", &auxSeconds)
	_, _ = fmt.Fscanf(r, "%d")
	//Searching errors process
	day, err[0] = strconv.Atoi(auxDay)
	month, err[1] = strconv.Atoi(auxMonth)
	year, err[2] = strconv.Atoi(auxYear)
	hour, err[3] = strconv.Atoi(auxHour)
	minutes, err[4] = strconv.Atoi(auxMinutes)
	seconds, err[5] = strconv.Atoi(auxSeconds)
	//Check if the date have errors
	for i := 0; i < len(err); i++ {
		if err[i] != nil {
			return nil, errors.New("Formato de fecha incorrecto.")
		}
	}
	//No error found then create the date
	parseMonth := time.Month(month)
	location, _ := time.LoadLocation("America/Argentina/Cordoba")
	auxDate := time.Date(year, parseMonth, day, hour, minutes, seconds, 0, location)
	auxUnixDate := auxDate.Unix()
	s := []byte(strconv.FormatInt(auxUnixDate, 10))
	unixDate := []byte(s)
	for i := len(unixDate); i < 10; i = len(unixDate) {
		unixDate = append([]byte{48}, unixDate...)
	}
	clearScreen()
	return unixDate, nil
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Println("#####################################")
	fmt.Println("___________Hamming/Huffman___________")
	fmt.Println("#####################################")
	fmt.Println()
}
