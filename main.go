package main

import (
	"bufio"
	"fmt"
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
		fmt.Println("6 - Salir")
		mainOp = 0
		_, _ = fmt.Fscanf(r, "%d", &mainOp)
		switch mainOp {
		case 1:
			Hamming()
		case 2:
			DeHamming()
		case 3:
			IntroduceErrors()
		case 5:
			seeSize()
		case 6:
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
