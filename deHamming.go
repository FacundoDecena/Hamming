package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func DeHamming() {
	var dhOp int
	r := bufio.NewReader(os.Stdin)
	dhContinue_ := true
	for dhContinue_ {
		clearScreen()
		fmt.Println("¿Que tipo de hamming ha sido aplicado?")
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
			preDeHamming7()
		case 5:
			dhContinue_ = false
		}
		_, _ = fmt.Fscanf(r, "%d")
	}
}

func preDeHamming7() {
	var fileName string
	var body []byte
	var err error
	r := bufio.NewReader(os.Stdin)

	clearScreen()
	fmt.Println("Ingrese el nombre del archivo .ha1")
	_, _ = fmt.Fscanf(r, "%s", &fileName)

	fileName += ".ha1"

	body, err = loadFile(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	decodedFile := deHamming7(body)
	fileName = strings.Replace(fileName, ".ha1", ".deh", -1)
	err = saveFile(fileName, decodedFile)
	if err != nil {
		fmt.Println(err)
	}

	_, _ = fmt.Scanf("%s")
	_, _ = fmt.Scanf("%s")
}

func deHamming7(file []byte) (ret []byte) {
	var encoded1stByte, encoded2ndByte, bitsToSpare byte
	two55 := exp(8) - 1 // 255
	var j byte
	j = 0
	for i := 0; i < len(file); i++ {
		if j == 0 {
			//Select the first 7-j bits
			encoded1stByte = file[i] & (two55 << (j + 1))
			//Save bits that does not belong to the hamming block
			bitsToSpare = file[i] & (exp(j+1) - 1)
			//Append decoded half to ret
			ret = append(ret, decode7(encoded1stByte))
			j++
			//Move bits to their place
			bitsToSpare = bitsToSpare<<8 - j
			//Select second hamming block
			encoded2ndByte = file[i+1] & (two55 << (j + 1))
			//Move the slice of block to its position
			encoded2ndByte = encoded2ndByte >> (j + 1)
			//Append bits to spare and the bits that belongs to the second hamming block
			encoded2ndByte = bitsToSpare & encoded2ndByte
			//Save bits that does not belong to the hamming block for the next iteration
			bitsToSpare = file[i+1] & (exp(j+1) - 1)
			//Append decoded half to ret
			ret = append(ret, decode7(encoded2ndByte))
			bitsToSpare = file[i] & (j + 1)
			j++
		} else {
			panic("Not implemented yet")
		}
		if j == 8 {
			j = 0
		}
	}
	return ret
}

func decode7(bait byte) (s byte) {
	c1 := (bait & uint8(1)) >> 7
	c2 := (bait & uint8(2)) >> 6
	d1 := (bait & uint8(4)) >> 5
	c3 := (bait & uint8(8)) >> 4
	d2 := (bait & uint8(16)) >> 3
	d3 := (bait & uint8(32)) >> 2
	d4 := (bait & uint8(64)) >> 1
	//Calculate sindrome using xor
	s1 := c1 ^ d1 ^ d2 ^ d4<<0
	s2 := c2 ^ d1 ^ d3 ^ d4<<1
	s3 := c3 ^ d2 ^ d3 ^ d4<<2

	s = s1 | s2 | s3

	if s != 0 {
		correct(bait)
	}

	d1 = d1 >> 4
	d2 = d2 >> 2
	d3 = d3 >> 1
	d4 = d4 >> 0

	return d1 | d2 | d3 | d4

}

func correct(wrong byte) (corrected byte) {
	return wrong
}

func exp(exponent byte) (ret byte) {
	ret = 1
	var i byte
	for i = 0; i < exponent; i++ {
		ret *= 2
	}
	return ret
}
