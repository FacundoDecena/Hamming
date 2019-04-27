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
	var body, testBody []byte
	var err error
	var start time.Time
	r := bufio.NewReader(os.Stdin)

	testBody = append(testBody, 205, 155, 54, 156, 208, 199, 165, 205, 155, 54, 156, 208, 199, 165)
	saveFile("test.ha1", testBody)
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
	//_, _ = fmt.Scanf("%s")
}

func deHamming7(file []byte) (ret []byte) {
	var encoded1stByte, encoded2ndByte, bitsToSpare, decoded1stByte, decoded2ndByte, decodedByte byte
	two55 := byte(exp(8)) - 1 // 255
	var j byte
	j = 0
	for i := 0; i < len(file); i += 2 {
		if j == 0 {
			//Select the first 7-j bits
			encoded1stByte = file[i] & (two55 << (j + 1))
			encoded1stByte >>= 1
			//Save bits that does not belong to the hamming block
			bitsToSpare = file[i] & (byte(exp(j+1)) - 1)
			//Append decoded half to decodedByte
			decoded1stByte = decode7(encoded1stByte) << 4
			j++
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
			bitsToSpare = file[i+1] & (byte(exp(j+1)) - 1)
			//Append 2nd decoded half to decodedByte
			decoded2ndByte = decode7(encoded2ndByte)
			decodedByte := decoded1stByte | decoded2ndByte
			//Append decodedByte to ret
			ret = append(ret, decodedByte)
			j++
			bitsToSpare = bitsToSpare << (8 - j)
		} else {
			//Select the first 7-j bits
			encoded1stByte = file[i] & (two55 << (j + 1))
			encoded1stByte >>= j

			encoded1stByte = bitsToSpare | encoded1stByte
			encoded1stByte >>= 1
			//Append decoded half to decodedByte
			decoded1stByte = decode7(encoded1stByte) << 4
			//Save bits that does not belong to the hamming block
			bitsToSpare = file[i] & (byte(exp(j+1)) - 1)
			j++
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
				bitsToSpare = file[i+1] & (byte(exp(j+1)) - 1)
				//Append 2nd decoded half to decodedByte
				decoded2ndByte = decode7(encoded2ndByte)
				decodedByte = decoded1stByte | decoded2ndByte
			}
			//Append decodedByte to ret
			ret = append(ret, decodedByte)
			j++
			bitsToSpare = bitsToSpare << (8 - j)

		}
		if j > 7 {
			j = 0
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
	s1 := c1 ^ d1 ^ d2 ^ d4<<0
	s2 := c2 ^ d1 ^ d3 ^ d4<<1
	s3 := c3 ^ d2 ^ d3 ^ d4<<2

	s = s1 | s2 | s3

	if s != 0 {
		correct(bait)
	}

	d1 = d1 << 4
	d2 = d2 << 2
	d3 = d3 << 1
	d4 = d4 << 0

	s = d1 | d2 | d3 | d4
	return s
}

func correct(wrong byte) (corrected byte) {
	return wrong
}

func exp(exponent byte) (ret int) {
	ret = 1
	var i byte
	for i = 0; i < exponent; i++ {
		ret *= 2
	}
	return ret
}
