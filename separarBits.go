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
		// var isEmpty bool
		fmt.Println("cuerpo original: ", body)
		fmt.Printf("cuerpo original: %08b\n", body)

		for body != nil {
			bytes, body = takeBits(26, body)

			/* haming algoritm */

			fmt.Printf("\nsaco los siguientes bits:\n bytes necesarios: %08b\n", bytes)
			fmt.Printf("cuerpo modificado %08b\n ", body)
			fmt.Println("cuerpo modificado: ", body)
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

			// for i := len(body); i <= cantByte; i++ {
			// 	body = append(body, uint8(0))
			// 	finish = true
			// }

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

			bait_aux := body[index]
			fmt.Printf("bait_aux: %08b    ", bait_aux)
			nextBait := body[index+1]
			fmt.Printf("nextByte: %08b\n", nextBait)

			nextBait = nextBait >> 6
			fmt.Printf("nextByte: %08b\n", nextBait)

			body[index] = bait_aux | nextBait
			body[index+1] = body[index+1] << 2

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
