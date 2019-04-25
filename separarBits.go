package main

import (
	"fmt"
	"math"
)

/*
func main() {
	fileName := "prueba.txt"
	var body []byte
	var bytes []byte
	TrashBits := 0

	body, err := loadFile(fileName)
	if err != nil {
		fmt.Println(err)
	} else {
		// var isEmpty bool
		fmt.Println("cuerpo original: ", body)
		fmt.Printf("cuerpo original: %08b\n", body)

		for body != nil {
			bytes, body, TrashBits = takeBits(26, body, TrashBits)

			//haming algoritm

			fmt.Printf("\nsaco los siguientes bits:\nbytes necesarios: %08b\n\n\n ", bytes)

		}
	}
}
*/
func takeBits(bits int, body []byte, NumberOfTrashBits int) ([]byte, []byte, int) {
	var cantByte int
	var cantBit int
	var arr_bit []byte
	var mask uint8
	var bait byte
	var finish bool

	cantByte = bits / 8 // amount bytes i need

	cantBit = bits - (cantByte * 8) // amount bits i need for the incomplete byte.

	if cantByte == 0 {
		finish = true

	} else {
		for index := 0; index < cantByte; index++ {
			if len(body) <= cantByte {
				finish = true
				body = append(body, uint8(0))
			}
			bitsToMove := 8 - NumberOfTrashBits // how many bits i need to shift.

			bait_aux := body[index]
			bait_aux = bait_aux << uint(NumberOfTrashBits) // shift to left how many bits i need to remove
			nextBait := body[index+1]
			nextBait = nextBait >> uint(bitsToMove)
			fmt.Printf("bait_aux : %08b  nextbait : %08b \n", bait_aux, nextBait)
			body[index] = bait_aux | nextBait // merge the bytes to pass the bits from the next byte to thisone.

			arr_bit = append(arr_bit, body[index])

		}
	}
	if finish == true {
		bait = body[cantByte] << uint(NumberOfTrashBits) // adjust the byte
		mask = doMask(cantBit)                           // make the mask by how many bits i need
		bait = bait & mask                               // make the byte
		fmt.Printf("voy a insertar %08b  mascara %08b\n", bait, mask)
		arr_bit = append(arr_bit, bait)
		body = nil
		body = nil
	} else {
		bait = body[cantByte] << uint(NumberOfTrashBits) // adjust the byte
		mask = doMask(cantBit)                           // make the mask by how many bits i need
		bait = bait & mask                               // make the byte
		fmt.Printf("voy a insertar %08b  mascara %08b\n", bait, mask)
		arr_bit = append(arr_bit, bait) // put the byte on the array.
		NumberOfTrashBits += cantBit
		body = body[cantByte:] // adjust the array

		// fmt.Printf("\nsaco los siguientes bits:\nbytes necesarios: %08b\n\n\n", arr_bit)
		// fmt.Printf("cuerpo %08b\n", body)
		fmt.Printf("numberoftrashbits: %d \n\n ", NumberOfTrashBits)

	}
	return arr_bit, body, NumberOfTrashBits
}

func doMask(bits int) uint8 {
	if bits > 8 {
		fmt.Printf("ERROR: WRONG MASK \n")

	} else if bits < 0 {
		fmt.Printf("ERROR: WRONG MASK \n")
	} else {
		val_mask := math.Pow(2, float64(bits)) - 1
		mask := uint8(val_mask) << uint(8-bits)
		return mask
	}
	return uint8(0)

}
