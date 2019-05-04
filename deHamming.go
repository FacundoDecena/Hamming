package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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
		case 2:
			preDeHamming(32)
		case 3:
			preDeHamming(1024)
		case 4:
			preDeHamming(32768)
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

func preDeHamming(size int) {
	var fileName string
	var body []byte
	var err error
	var start time.Time
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	var format string
	switch size {
	case 32:
		format = "2"
	case 1028:
		format = "3"
	case 32768:
		format = "4"

	}
	fmt.Println("Ingrese el nombre del archivo .ha" + format)
	_, _ = fmt.Fscanf(r, "%s", &fileName)

	fileName += ".ha" + format

	body, err = loadFile(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}
	start = time.Now()
	decodedFile := callDecode(size, body)
	fileName = strings.Replace(fileName, ".ha"+format, ".deh", -1)
	err = saveFile(fileName, decodedFile)
	if err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	log.Printf("\nDeHamming took %s", elapsed)
	_, _ = fmt.Scanf("%s")
}
func deHamming7(file []byte) (ret []byte) {
	var encoded1stByte, encoded2ndByte, bitsToSpare, decoded1stByte, decoded2ndByte, decodedByte byte
	bitsToSpare = 0
	two55 := byte(exp(8)) - 1 // 255
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
		bitsToSpare = file[i] & (byte(exp(j+1)) - 1)
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
			bitsToSpare = file[i+1] & (byte(exp(j+1)) - 1)
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
	mistake := bait & byte(exp(7-syndrome))
	//If it is 0, the only way to know its position is using the same power of 2
	if mistake == 0 {
		mistake = byte(exp(7 - syndrome))
	} else { //If the bit is 1 it has to be 0
		mistake = 0
	}
	//wom comes from Without Mistake, which is bait with 0 in the position of the mistake
	wom := bait & (255 - byte(exp(7-syndrome)))

	corrected = wom | mistake

	return corrected
}

func exp(exponent byte) (ret int) {
	ret = 1
	var i byte
	for i = 0; i < exponent; i++ {
		ret *= 2
	}
	return ret
}

func callDecode(size int, input []byte) []byte {
	var decodedFile []byte
	_, _, controlBitsQuantity := initialCase(size)
	blockSize := size / 8
	il := 0
	sl := blockSize
	for i := 0; i < len(input); i += blockSize {
		aux := decode(size, input[il:sl], controlBitsQuantity)
		for j := 0; j < len(aux); j++ {
			decodedFile = append(decodedFile, aux[j])
		}
		il += blockSize
		sl += blockSize
	}
	return compress(decodedFile)
}

func decode(size int, input []byte, controlBitsQuantity int) []byte {
	decoded := make([]byte, int(math.Ceil((float64(size)-float64(controlBitsQuantity))/8)))
	decodedPosition := 0
	decodedNumberOfByte := 0
	for i := 0; i < controlBitsQuantity-1; i++ {
		il := expInt(i) - 1
		sl := expInt(i+1) - 1
		for j := il + 1; j < sl; j++ {
			position, numberOfByte := byteNumberDeHamming(j)
			dataBit := takeBit(input[numberOfByte], 7-position, 7-decodedPosition)
			decoded[decodedNumberOfByte] = decoded[decodedNumberOfByte] | dataBit
			decodedPosition++
			if decodedPosition > 7 {
				decodedPosition = 0
				decodedNumberOfByte++
			}
		}
	}
	return decoded
}

func byteNumberDeHamming(position int) (int, int) {
	place := position % 8
	byteNumber := position / 8
	return place, byteNumber
}

func compress(input []byte) []byte {
	have := make([]int, len(input)/4)
	need := make([]int, len(input)/4)
	var compressed []byte //:= make([]byte, int(math.Ceil(float64(len(input))*3.25/float64(4))))
	for i := 0; i < len(input)/4-1; i++ {
		have[i] = 26
		need[i] = 6
	}
	position := 2
	for i := 0; i < len(input)/4-1; i++ {
		if need[i]%8 == 0 {
			position = 2
			for j := 0; j < have[i]/8; j++ {
				compressed = append(compressed, input[i*4+j])
			}
			for j := 0; j < need[i]/8; j++ {
				compressed /*[i*4+have[i]/8+j]*/ = append(compressed, input[(i+1)*4+j])
			}
			need[i+1] += need[i]
			have[i+1] -= need[i]
			index := 0
			for j := need[i] / 8; j < 4; j++ {
				input[(i+1)*4+index] = input[(i+1)*4+j]
				index++
			}
			for j := index; j < 4; j++ {
				input[(i+1)*4+j] = byte(0)
			}
		} else if need[i] > have[i+1] && i != len(input)/4-2 {
			aux1 := takeBitsDeHamming(have[i+1], input[(i+1)*4:(i+1)*4+5], position)
			aux2 := takeBitsDeHamming(need[i]-have[i+1], input[(i+2)*4:(i+2)*4+5], (position+have[i+1])%8)
			conflictByte1 := int(have[i] / 8)
			for j := 0; j < conflictByte1; j++ {
				compressed /*[i*4+j]*/ = append(compressed, input[i*4+j])
			}
			compressed /*[i*4+conflictByte1]*/ = append(compressed, input[i*4+conflictByte1]|aux1[0])
			for j := conflictByte1 + 1; j < len(aux1)-1; j++ {
				compressed /*[i*4+j]*/ = append(compressed, aux1[j])
			}
			conflictByte2 := int(have[i+1] / 8)
			compressed /*[i*4+len(aux1)-1]*/ = append(compressed, aux1[conflictByte2]|aux2[0])
			for j := 1; j < len(aux2); j++ {
				compressed /*[i*4+len(aux1)+j]*/ = append(compressed, aux2[j])
			}
			aux := ajustBytes(input[(i+2)*4:(i+2)*4+5], need[i]-have[i+1])
			for j := 0; j < len(aux); j++ {
				input[(i+2)*4+j] = aux[j]
			}
			have[i+2] -= need[i] - have[i+1]
			need[i+2] += need[i] - have[i+1]
			i++
			position = (position + 4) % 8
		} else {
			aux := takeBitsDeHamming(need[i], input[(i+1)*4:(i+1)*4+5], position)
			conflictByte := int(have[i] / 8)
			for j := 0; j < conflictByte; j++ {
				compressed /*[i*4+j]*/ = append(compressed, input[i*4+j])
			}
			compressed /*[i*4+conflictByte]*/ = append(compressed, input[i*4+conflictByte]|aux[0])
			for j := 1; j < len(aux); j++ {
				compressed /*[i*4+j]*/ = append(compressed, aux[j])
			}
			need[i+1] += need[i]
			have[i+1] -= need[i]
			aux = ajustBytes(input[(i+1)*4:(i+1)*4+4+1], need[i])
			for j := 0; j < len(aux); j++ {
				input[(i+1)*4+j] = aux[j]
			}
			if need[i] == 26 {
				i++
				position = (position + 4) % 8
			} else {
				position += 2
			}
		}
	}
	for i := 0; i < 4; i++ {
		compressed = append(compressed, input[len(input)-4+i])
	}
	return compressed[:]
}

func takeBitsDeHamming(bits int, input []byte, initialPosition int) []byte {
	aux := input[len(input)-1]
	input[len(input)-1] = byte(0)
	bytesQuantity := int(math.Ceil(float64(bits+initialPosition) / float64(8)))
	ret := make([]byte, bytesQuantity)
	if initialPosition == 0 {
		for i := 0; i < bytesQuantity; i++ {
			ret[i] = input[i]
		}
		if bits%8 != 0 {
			ret[bytesQuantity-1] &= doMask(bits % 8)
		}
	} else {
		garbage := byte(0)
		for i := 0; i < bytesQuantity; i++ {
			ret[i] = garbage | ((doMask(8-initialPosition) & input[i]) >> byte(initialPosition))
			garbage = ((doMask(initialPosition) >> byte(8-initialPosition)) & input[i]) << byte(8-initialPosition)
		}
		mask := (bits%8 + initialPosition) % 8
		if mask == 0 {
			mask = 8
		}
		ret[bytesQuantity-1] &= doMask(mask)
	}
	input[len(input)-1] = aux
	return ret
}

func ajustBytes(input []byte, begin int) []byte {
	position := begin % 8
	numberOfByte := byteNumber(begin, 4)
	bytesQuantity := 4
	ret := make([]byte, bytesQuantity)
	ret[0] = ((doMask(8-position) >> byte(position)) & input[numberOfByte]) << byte(position)
	aux := takeBitsDeHamming((len(input)-numberOfByte-2)*8, input[numberOfByte+1:], 8-position)
	ret[0] |= aux[0]
	for i := 1; i < len(aux); i++ {
		ret[i] = aux[i]
	}
	return ret
}
