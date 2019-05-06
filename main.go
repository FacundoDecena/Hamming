package main

import (
	"bufio"
	"fmt"
	//Para que le ande a ustedes
	//"github.com/FacundoDecena/Hamming/Huffman"
	//Para que me ande a mi
	"Hamming/Huffman"
	"os"
	"os/exec"
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

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Println("#####################################")
	fmt.Println("_______________HAMMING_______________")
	fmt.Println("#####################################")
	fmt.Println()
}
