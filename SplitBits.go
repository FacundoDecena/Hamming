package main

import (
	"fmt"
	"math"
)

func callTakeBits(body []byte) (Array [][]byte) {
	var TrashBits int
	var bytes []byte

	TrashBits = 0

	for body != nil {
		bytes, body, TrashBits = takeBits(1013, body, TrashBits)

		Array = append(Array, bytes)
	}
	return Array

}

func takeBits(bits int, body []byte, NumberOfTrashBits int) ([]byte, []byte, int) {

	var cantByte int
	var cantBit int
	var arr_bit []byte
	var mask uint8
	var bait byte
	var finish bool
	var bitsToMove int

	// initialise variables //

	cantByte = bits / 8                // amount bytes i need.
	bitsToMove = 8 - NumberOfTrashBits // how many bits i need to shift.

	cantBit = bits - (cantByte * 8) // amount bits i need for the incomplete byte.

	if bits >= 8 {

		for index := 0; index < cantByte; index++ {
			if len(body) <= cantByte {
				body = append(body, uint8(0))
				finish = true
			}

			// adjust //

			bait_aux := body[index]
			bait_aux = bait_aux << uint(NumberOfTrashBits) // shift to left how many bits i need to remove
			nextBait := body[index+1]
			nextBait = nextBait >> uint(bitsToMove)
			fmt.Printf("bait_aux  %d : %08b d nextbait : %08b \n", index, bait_aux, nextBait)
			body[index] = bait_aux | nextBait // merge the bytes to pass the bits from the next byte to thisone.

			arr_bit = append(arr_bit, body[index]) // put it on the array.

		}
		if finish == true {
			body = nil
		} else {

			bait = body[cantByte] << uint(NumberOfTrashBits) // adjust the byte

			mask = doMask(cantBit) // make the mask by how many bits i need
			bait = bait & mask     // make the byte
			fmt.Printf("voy a insertar %08b  mascara %08b\n", bait, mask)
			arr_bit = append(arr_bit, bait) // put the byte on the array.
			NumberOfTrashBits += cantBit
			body = body[cantByte:] // adjust the array

			fmt.Printf("\nsaco los siguientes bits:\nbytes necesarios: %08b\n\n\n", arr_bit)
			fmt.Printf("numberoftrashbits: %d \n\n ", NumberOfTrashBits)

		}
		return arr_bit, body, NumberOfTrashBits
	} else {
		fmt.Printf("This function is not available for values less than 8 bits.")
		return nil, nil, 0
	}
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