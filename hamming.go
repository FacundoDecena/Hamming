package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

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
			preHamming7()
		case 5:
			dhContinue_ = false
		}
		_, _ = fmt.Fscanf(r, "%d")
	}

}

func preHamming7() {
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
		start = time.Now()
		encodedBody = hamming7(body)
		fileName = strings.Replace(fileName, ".txt", ".ha1", -1)
		err = saveFile(fileName, encodedBody)
		if err != nil {
			fmt.Println(err)
		}

	}
	elapsed := time.Since(start)
	log.Printf("\nHamming7 took %s", elapsed)
	_, _ = fmt.Fscanf(r, "%s")
	_, _ = fmt.Fscanf(r, "%s")
}

// Receives a byte slice, returns it encoded
func hamming7(file []byte) []byte {
	//Mask that shows first bits
	maskFirst := 240
	//Mask that shows last bits
	maskLast := 15
	//Encoded byte slice
	var ret []byte
	//Var j indicates the index of block attached
	j := 0
	//x is a complete byte, y is an incomplete byte
	var x, y byte

	for i := 0; i < len(file); i++ {
		//Add mask
		hammingMaskedFirst := make([]byte, 1)
		hammingMaskedFirst[0] = file[i] & uint8(maskFirst)
		//Move to last positions
		hammingMaskedFirst[0] /= 16
		//Add mask
		hammingMaskedLast := make([]byte, 1)
		hammingMaskedLast[0] = file[i] & uint8(maskLast)
		//Put extra bit on the right
		maskedFirst := encode(8, hammingMaskedFirst)[0] << 1
		//Put extra bit on the right
		maskedLast := encode(8, hammingMaskedLast)[0] << 1
		//If j==0 maskedFirst can go directly to ret, otherwise has to wait next byte to complete
		if j == 0 {
			maskedFirst, maskedLast = compress7(maskedFirst, maskedLast, j)
			ret = append(ret, maskedFirst)
			y = maskedLast
		} else if j < 8 {
			//now y is complete, then is saved in x and can be append to ret
			x, y = compress7(y, maskedFirst, j)
			ret = append(ret, x)
			//we can complete y with maskedLast first j bits, after that y meets the conditions to be x
			x, y = compress7(y, maskedLast, j)
			ret = append(ret, x)
		} else {
			//Restart the process
			j = 0
		}
		j++
	}
	return ret
}

//Joins x and y, returns y moved to the left
func compress7(x byte, y byte, index int) (byte, byte) {
	//My implementation for **
	exp := 1
	i := uint8(7 - index)
	for ; i > 0; i-- {
		exp *= 2
	}
	//Mask the bits we need
	y1 := y & uint8(exp) >> uint8(7-index)
	//Join x with the bits above
	x = x | y1
	//Move y index+1 places to the left
	y = y << uint8(index+1)

	return x, y
}

//Size should be: 8 for hamming7, 32 for hamming 32, 1024 for hamming 1024 and 32768 for hamming 32768
func encode(size int, input []byte) []byte {
	encoded := make([]byte, int(size/8))
	var position int
	var numberOfByte int
	var controlBitsQuantity int
	//Set the initial position where is the first information bit in the array passed by parameter depending of what hamming will be apply
	switch size {
	case 8:
		position = 0
		numberOfByte = 0
		controlBitsQuantity = 2
	case 32:
		position = 6
		numberOfByte = 3
		controlBitsQuantity = 4
	case 1024:
		position = 3
		numberOfByte = 126
		controlBitsQuantity = 9
	case 32768:
		position = 0
		numberOfByte = 4095
		controlBitsQuantity = 14
	default:
		return nil
	}
	//Data bits accommodate process
	for i := 0; i < controlBitsQuantity+1; i++ {
		il := expInt(i) - 1
		sl := expInt(i+1) - 1
		for j := il + 1; j < sl; j++ {
			dataBit := takeBit(input[numberOfByte], position, int(j%8))
			x := byteNumber(int(j), size/8)
			encoded[x] = encoded[x] | dataBit
			position++
			if position > 7 {
				numberOfByte--
				position = 0
			}
		}
	}
	//Control bits calculus process
	for i := 0; i < controlBitsQuantity+1; i++ {
		parity := byte(0)
		for j := expInt(i) - 1; j < size; j += expInt(i + 1) {
			for k := 0; k < expInt(i); k++ {
				parity = parity ^ takeBit(encoded[byteNumber(j+int(k), size/8)], int((j+k)%8), 0)
			}
		}
		if takeBit(parity, 0, 0) != 0 {
			x := byteNumber(int(expInt(i)-1), size/8)
			encoded[x] = encoded[x] | takeBit(1, 0, int(expInt(i)-1)%8)
		}
	}
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
