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

var tiempoControlCalculus time.Duration
var tiempoAccomodateBits time.Duration
var tiempoAccomodateArray time.Duration
var tiempoTakeBits time.Duration

func Hamming() {
	var dhOp int
	r := bufio.NewReader(os.Stdin)
	dhContinue_ := true
	for dhContinue_ {
		clearScreen()
		fmt.Println("¿Que tipo de hamming quiere aplicar?")
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
	fmt.Println("Ingrese el nombre del archivo")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Since golang does not show the time a program runs...
	start := time.Now()

	body, err := loadFile(fileName)

	if err != nil {
		fmt.Println(err)
	} else {
		//Start the timer
		//start = time.Now()
		switch size {
		case 7:
			encodedBody = hamming7(body)
		case 32:
			encodedBody = hamming(size, body)
		case 1024:
			encodedBody = hamming(size, body)
		case 32768:
			encodedBody = hamming(size, body)

		}

		fileName = strings.Replace(fileName, ".txt", ".ha1", -1)
		err = saveFile(fileName, encodedBody)
		if err != nil {
			fmt.Println(err)
		}

	}
	elapsed := time.Since(start)
	log.Printf("\nHamming took %s", elapsed)
	log.Printf("\nCalculoBits took %s", tiempoControlCalculus)
	log.Printf("\nAcomodarArreglo took %s", tiempoAccomodateArray)
	log.Printf("\nAcomodarBits took %s", tiempoAccomodateBits)
	log.Printf("\nTomarBits took %s", tiempoTakeBits)
	_, _ = fmt.Fscanf(r, "%s")
	_, _ = fmt.Fscanf(r, "%s")
}

// Receives a byte slice, returns it encoded
func hamming7(file []byte) []byte {
	//Mask that shows first bits
	mask1 := 240
	//Mask that shows last bits
	mask2 := 15
	entryLength := len(file)
	//Number that I use so that the size of the array is a multiple of 8, thus making compression simpler
	module := 0
	if 2*entryLength%8 != 0 {
		module = 8 - 2*entryLength%8
	}
	auxLength := 2*entryLength + module
	finalLength := int(math.Ceil(float64(entryLength) * 1.75))
	var auxArray = make([]byte, auxLength)
	//Applies the hamming encode to each byte of the file
	for i := 0; i < entryLength; i++ {
		var firstBits, lastBits byte
		firstBits = (file[i] & uint8(mask1)) >> 4
		lastBits = file[i] & uint8(mask2)
		auxArray[2*i] = encode7(firstBits)
		auxArray[2*i+1] = encode7(lastBits)
	}
	j := 0
	ret := make([]byte, auxLength)
	//Compress the array
	for i := 0; i < auxLength; i += 8 {
		sevenBlock := compressBlock(auxArray[i : i+8])
		ret[j] = sevenBlock[0]
		ret[j+1] = sevenBlock[1]
		ret[j+2] = sevenBlock[2]
		ret[j+3] = sevenBlock[3]
		ret[j+4] = sevenBlock[4]
		ret[j+5] = sevenBlock[5]
		ret[j+6] = sevenBlock[6]
		j += 7
	}
	return ret[0:finalLength]
}

func encode7(bait byte) byte {
	//Get bits from position in brackets and send it to the left
	d4 := bait & uint8(1)
	d3 := (bait & uint8(2)) >> 1
	d2 := (bait & uint8(4)) >> 2
	d1 := (bait & uint8(8)) >> 3
	//Calculate controls using xor
	c1 := d1 ^ d2 ^ d4
	c2 := d1 ^ d3 ^ d4
	c3 := d2 ^ d3 ^ d4
	//set variables in their position
	c1 <<= 7
	c2 <<= 6
	d1 <<= 5
	c3 <<= 4
	d2 <<= 3
	d3 <<= 2
	d4 <<= 1
	return d4 | d3 | d2 | c3 | d1 | c2 | c1
}

func compressBlock(bp []byte) [7]byte {
	var ba [7]byte
	ba[0] = bp[0]
	ba[0] = ba[0] | ((bp[1] & 128) >> 7)
	ba[1] = bp[1] << 1
	ba[1] = ba[1] | ((bp[2] & 192) >> 6)
	ba[2] = bp[2] << 2
	ba[2] = ba[2] | ((bp[3] & 224) >> 5)
	ba[3] = bp[3] << 3
	ba[3] = ba[3] | ((bp[4] & 240) >> 4)
	ba[4] = bp[4] << 4
	ba[4] = ba[4] | ((bp[5] & 248) >> 3)
	ba[5] = bp[5] << 5
	ba[5] = ba[5] | ((bp[6] & 252) >> 2)
	ba[6] = bp[6] << 6
	ba[6] = ba[6] | (bp[7] >> 1)
	return ba
}

func hamming(size int, file []byte) []byte {
	var ret []byte
	switch size {
	case 32:
		start := time.Now()
		x := callTakeBits(26, file)
		tiempoTakeBits += time.Since(start)
		ret = callEncode(size, x)

	case 1024:
		x := callTakeBits(1013, file)
		ret = callEncode(size, x)
	case 32768:
		x := callTakeBits(32752, file)
		ret = callEncode(size, x)
	}
	return ret
}

func callEncode(size int, inputFile [][]byte) (outPut []byte) {
	position, numberOfByte, controlBitsQuantity := initialCase(size)
	var aux [][]byte
	for i := 0; i < len(inputFile); i++ {
		aux = append(aux, encode(size, inputFile[i], position, numberOfByte, controlBitsQuantity))
		start := time.Now()
		for j := 0; j < len(aux[i]); j += len(aux[i]) {
			outPut = append(outPut, aux[i][j])
		}
		tiempoAccomodateArray += time.Since(start)
		fmt.Printf("%d %d %fcompletado\n", i, len(inputFile), int(((i+1)/len(inputFile))*100))
	}
	return outPut
}

//Size should be: 8 for hamming7, 32 for hamming 32, 1024 for hamming 1024 and 32768 for hamming 32768
func encode(size int, input []byte, position int, numberOfByte int, controlBitsQuantity int) []byte {
	encoded := make([]byte, int(size/8))
	start1 := time.Now()
	//Data bits accommodate process
	for i := controlBitsQuantity - 1; i > 0; i-- {
		sl := expInt(i) - 1
		il := expInt(i-1) - 1
		for j := sl - 1; j > il; j-- {
			dataBit := takeBit(input[numberOfByte], position, 7-int(j%8))
			x := byteNumber(int(j), size/8)
			encoded[x] = encoded[x] | dataBit
			position++
			if position > 7 {
				numberOfByte--
				position = 0
			}
		}
	}
	tiempoAccomodateBits += time.Since(start1)
	start2 := time.Now()
	//Control bits calculus process
	for i := 0; i < controlBitsQuantity-1; i++ {
		parity := byte(0)
		for j := expInt(i) - 1; j < size; j += expInt(i + 1) {
			for k := 0; k < expInt(i); k++ {
				parity ^= takeBit(encoded[byteNumber(j+int(k), size/8)], 7-(int((j+k)%8)), 0)
			}
		}
		x := byteNumber(int(expInt(i)-1), size/8)
		encoded[x] = encoded[x] | takeBit(1, 0, 7-(int(expInt(i)-1)%8))
	}
	tiempoControlCalculus += time.Since(start2)
	return encoded
}

//Apply a mask to a source byte to get the bit in the initial position and shifter it to the final position.
func takeBit(source byte, initialPosition int, finalPosition int) byte {
	var result = source & byte(expInt(initialPosition))
	var shift = finalPosition - initialPosition
	if shift == 0 {
		return result
	} else if shift > 0 {
		return result << uint(shift)
	} else {
		shift *= -1
		return result >> uint(shift)
	}
}

//From a bit position in a block of bytes of size bytesQuantity, returns the number of the byte which the bit belongs
func byteNumber(bitPosition int, bytesQuantity int) int {
	il := 0
	sl := 7
	for i := bytesQuantity - 1; i > 0; i-- {
		if bitPosition >= il && bitPosition <= sl {
			return i
		} else {
			il += 8
			sl += 8
		}
	}
	return 0
}

func expInt(exponent int) int {
	//My implementation for **
	var result = 1
	for i := 0; i < exponent; i++ {
		result *= 2
	}
	return result
}

func initialCase(size int) (position int, numberOfByte int, controlBitsQuantity int) {
	//Set the initial position where is the first information bit in the array passed by parameter depending of what hamming will be apply
	switch size {
	case 32:
		position = 6
		numberOfByte = 3
		controlBitsQuantity = 6
	case 1024:
		position = 3
		numberOfByte = 126
		controlBitsQuantity = 11
	case 32768:
		position = 0
		numberOfByte = 4095
		controlBitsQuantity = 16
	}
	return position, numberOfByte, controlBitsQuantity
}
