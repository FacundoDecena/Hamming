package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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
	var start time.Time
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
	start = time.Now()
	decodedFile := deHamming7(body)
	fileName = strings.Replace(fileName, ".ha1", ".deh", -1)
	err = saveFile(fileName, decodedFile)
	if err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	log.Printf("\nDeHamming7 took %s", elapsed)
	_, _ = fmt.Scanf("%s")
}

func deHamming7(file []byte) (ret []byte) {
	var encoded1stByte, encoded2ndByte, bitsToSpare, decoded1stByte, decoded2ndByte, decodedByte byte
	bitsToSpare = 0
	two55 := exp(8) - 1 // 255
	var j byte
	j = 0
	for i := 0; i < len(file); i += 2 {
		//Select the first 7-j bits
		encoded1stByte = file[i] & (two55 << (j + 1))
		//Move them to theirs position
		encoded1stByte >>= j
		//Join the pieces
		encoded1stByte = bitsToSpare | encoded1stByte
		//Move the leftover bit to the left
		encoded1stByte >>= 1
		//Append decoded half to decodedByte
		decoded1stByte = decode7(encoded1stByte) << 4
		//Save bits that does not belong to the hamming block
		bitsToSpare = file[i] & (exp(j+1) - 1)
		j++
		if j%7 == 0 && i > 0 {
			i--
			decodedByte = decoded1stByte | bitsToSpare
		}
		if i+1 == len(file) {
			decodedByte = decoded1stByte | bitsToSpare
		} else {
			//Move bits to their place
			bitsToSpare = bitsToSpare << (8 - j)
			//Select second hamming block
			encoded2ndByte = file[i+1] & (two55 << (j + 1))
			//Move the slice of block to its position
			encoded2ndByte = encoded2ndByte >> (j)
			//Append bits to spare and the bits that belongs to the second hamming block
			encoded2ndByte = bitsToSpare | encoded2ndByte
			encoded2ndByte >>= 1
			//Save bits that does not belong to the hamming block for the next iteration
			bitsToSpare = file[i+1] & (exp(j+1) - 1)
			//Append 2nd decoded half to decodedByte
			decoded2ndByte = decode7(encoded2ndByte)
			decodedByte = decoded1stByte | decoded2ndByte
		}
		//Append decodedByte to ret
		ret = append(ret, decodedByte)
		j++
		bitsToSpare = bitsToSpare << (8 - j)
		if j > 7 {
			j = 0
			bitsToSpare = 0
		}
	}
	return ret
}

func decode7(bait byte) (s byte) {
	c1 := (bait & uint8(64)) >> 6
	c2 := (bait & uint8(32)) >> 5
	d1 := (bait & uint8(16)) >> 4
	c3 := (bait & uint8(8)) >> 3
	d2 := (bait & uint8(4)) >> 2
	d3 := (bait & uint8(2)) >> 1
	d4 := (bait & uint8(1)) >> 0
	//Calculate sindrome using xor
	var s1, s2, s3 byte

	s1 = (c1 ^ d1 ^ d2 ^ d4) << 0
	s2 = (c2 ^ d1 ^ d3 ^ d4) << 1
	s3 = (c3 ^ d2 ^ d3 ^ d4) << 2

	s = s1 | s2 | s3

	if s != 0 {
		bait = correct(bait, s)
	}

	d1 = (bait & uint8(16)) >> 4
	d2 = (bait & uint8(4)) >> 2
	d3 = (bait & uint8(2)) >> 1
	d4 = (bait & uint8(1)) >> 0

	d1 = d1 << 3
	d2 = d2 << 2
	d3 = d3 << 1
	d4 = d4 << 0

	s = d1 | d2 | d3 | d4
	return s
}

func correct(bait byte, syndrome byte) (corrected byte) {
	//Get the wrong bit
	mistake := bait & exp(7-syndrome)
	//If it is 0, the only way to know its position is using the same power of 2
	if mistake == 0 {
		mistake = exp(7 - syndrome)
	} else { //If the bit is 1 it has to be 0
		mistake = 0
	}
	//wom comes from Without Mistake, which is bait with 0 in the position of the mistake
	wom := bait & (255 - exp(7-syndrome))

	corrected = wom | mistake

	return corrected
}

func exp(exponent byte) (ret byte) {
	ret = 1
	var i byte
	for i = 0; i < exponent; i++ {
		ret *= 2
	}
	return ret
}
