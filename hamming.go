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
		maskedFirst := file[i] & uint8(maskFirst)
		//Move to last positions
		maskedFirst /= 16
		//Add mask
		maskedLast := file[i] & uint8(maskLast)
		//Put extra bit on the right
		maskedFirst = encode7(maskedFirst)
		//Put extra bit on the right
		maskedLast = encode7(maskedLast)
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

func encode7(bait byte) byte {
	//Get bits from position in brackets and send it to the left
	d1 := (bait & uint8(1)) >> 0
	d2 := (bait & uint8(2)) >> 1
	d3 := (bait & uint8(4)) >> 2
	d4 := (bait & uint8(8)) >> 3
	//Calculate controls using xor
	c1 := d1 ^ d2 ^ d4
	c2 := d1 ^ d3 ^ d4
	c3 := d2 ^ d3 ^ d4
	//set variables in their position
	c1 = c1 << 1
	c2 = c2 << 2
	d1 = d1 << 3
	c3 = c3 << 4
	d2 = d2 << 5
	d3 = d3 << 6
	d4 = d4 << 7

	return d4 | d3 | d2 | c3 | d1 | c2 | c1
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

func encode32(input [4]byte)[4] byte{
	var encoded [4]byte
	//The 6th bit of input[3] is the less significant bit where i begin to accommodate the bits in encoded
	var position = 6
	var numberOfByte = 3
	//Data bits accommodate process
	for i:=0; i<5; i++ {
		il := exp(i) - 1
		sl := exp(i + 1)  - 1
		for j := il + 1; j < sl; j++ {
			var dataBit = takeBit(input[numberOfByte], position, int(j%8))
			var x = byteNumber(int(j),4)
			encoded[x] = encoded[x] | dataBit
			position++
			if position > 7 {
				numberOfByte--
				position = 0
			}
		}
	}
	//Control bits calculus process
	for i:=0;i<5;i++{
		var parity = byte(0)
		for j := exp(i)-1; j<32; j+= exp(i+1){
			for k := 0; k < exp(i); k++{
				parity = parity ^ takeBit(encoded[byteNumber(j+int(k), 4)], int((j+k)%8), 0)
			}
		}
		if takeBit(parity,0,0) != 0{
			x := byteNumber(int(exp(i)-1),4)
			encoded[x] = encoded[x] | takeBit(1,0,int(exp(i)-1)%8)
		}
	}
	return encoded
}

//Apply a mask to a source byte to get the bit in the initial position and shifter it to the final position.
func takeBit(source byte,initialPosition int,finalPosition int) byte{
	var result = source & byte(exp(initialPosition))
	var shift = finalPosition - initialPosition
	if shift == 0{
		return result
	}else if shift > 0{
		return result << uint(shift)
	}else{
		shift*=-1
		return result >> uint(shift)
	}
}
``
//From a bit position in a block of bytes of size bytesQuantity, returns the number of the byte which the bit belongs
func byteNumber(bitPosition int,bytesQuantity int) int{
	il := 0
	sl := 7
	for i:= bytesQuantity-1;i>0;i--{
		if bitPosition>=il && bitPosition <=sl{
			return i
		}else{
			il +=8
			sl +=8
		}
	}
	return 0
}
