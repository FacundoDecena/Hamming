package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	var fileName string
	var encodedBody []byte

	fmt.Println("Ingrese el nombre del archivo")
	_, _ = fmt.Scanf("%s", &fileName)

	start := time.Now()

	body, err := loadFile(fileName)

	if err != nil {
		fmt.Println(err)
	} else {
		start = time.Now()
		encodedBody = hamming7(body)

		err = saveFile("encoded"+fileName, encodedBody)
		if err != nil {
			fmt.Println(err)
		}

	}
	elapsed := time.Since(start)
	log.Printf("\nHamming7 took %s", elapsed)

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
		maskedFirst = encode(maskedFirst) << 1
		//Put extra bit on the right
		maskedLast = encode(maskedLast) << 1
		//If j==0 maskedFirst can go directly to ret, otherwise has to wait next byte to complete
		if j == 0 {
			maskedFirst, maskedLast = compress(maskedFirst, maskedLast, j)
			ret = append(ret, maskedFirst)
			y = maskedLast
		} else if j < 8 {
			x, y = compress(y, maskedFirst, j)
			ret = append(ret, x)
			x = y
			x, y = compress(y, maskedLast, j)
			ret = append(ret, x)
		} else {
			j = 0
		}
		j++
	}
	return ret
}

func encode(bait byte) byte {
	d4 := (bait & uint8(1)) >> 0
	d3 := (bait & uint8(2)) >> 1
	d2 := (bait & uint8(4)) >> 2
	d1 := (bait & uint8(8)) >> 3

	c1 := d1 ^ d2 ^ d4
	c2 := d1 ^ d3 ^ d4
	c3 := d2 ^ d3 ^ d4
	c2 = c2 << 1
	d1 = d1 << 2
	c3 = c3 << 3
	d2 = d2 << 4
	d3 = d3 << 5
	d4 = d4 << 6

	return d4 | d3 | d2 | c3 | d1 | c2 | c1
}

func compress(x byte, y byte, index int) (byte, byte) {
	exp := 1
	i := uint8(7 - index)
	for ; i > 0; i-- {
		exp *= 2
	}
	y1 := y & uint8(exp) >> uint8(7-index)
	x = x | y1
	y = y << uint8(index+1)

	return x, y
}

func loadFile(fileName string) ([]byte, error) {
	var err error
	var body []byte
	body, err = ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func saveFile(fileName string, body []byte) error {

	return ioutil.WriteFile(fileName, body, 0600)

}
