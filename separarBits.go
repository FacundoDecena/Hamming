package main

import (
	"fmt"
	"math"
)

func main() {
	fileName := "prueba.txt"
	var body []byte
	var bytes []byte

	body, err := loadFile(fileName)
	if err != nil {
		fmt.Println(err)
	} else {

		for body != nil {
			bytes, body = takeBits(26, body)

			/* haming algoritm */
			fmt.Printf("%08b", bytes)

		}

	}
}

func takeBits(bits int, body []byte) ([]byte, []byte) {
	var cantByte int
	var cantBit int
	var arr_bit []byte
	var mask uint8
	var bait byte
	var finish bool

	cantByte = bits / 8             // amount bytes i need
	cantBit = bits - (cantByte * 8) // amount bits i need for the incomplete byte.
	for index := 0; index < cantByte; index++ {
		if len(body) <= cantByte {
			finish = true
			body = append(body, uint8(0))

		}
		arr_bit = append(arr_bit, body[index])
	}
	if finish == true {
		arr_bit = append(arr_bit, body[cantByte])
		body = nil
	} else {

		mask = doMask(cantBit)               // make the mask by how many bits i need
		bait = body[cantByte] & mask         // make the byte
		arr_bit = append(arr_bit, bait)      // put the byte on the array.
		body[cantByte] = body[cantByte] << 2 // adjust the byte.
		body = body[cantByte:]               // adjust the array

		// adjust the array of bytes.

		for index := 0; index < len(body)-1; index++ {

			bait_aux := body[index] // take the byte

			nextBait := body[index+1] // take the next byte

			nextBait = nextBait >> 6 // move the bits i want to the right

			body[index] = bait_aux | nextBait // merge the bits of both bytes

			body[index+1] = body[index+1] << 2 // adjust the nextByte.

		}

	}
	return arr_bit, body
}

func doMask(bits int) uint8 {
	if bits > 8 {
		fmt.Printf("ERROR: WRONG MASK \n")

	} else if bits < 0 {
		fmt.Printf("ERROR: WRONG MASK \n")
	} else {
		val_mask := math.Pow(2, float64(bits)) - 1
		mask := uint8(val_mask) << uint((8 - bits))
		return mask
	}
	return uint8(0)

}
