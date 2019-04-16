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

func exp(exponent int) int {
	//My implementation for **
	var result = 1
	for i := 0; i < exponent; i++ {
		result *= result
	}
	return result
}
