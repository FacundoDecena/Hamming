package HammingCodification

import (
	"math"
	"strconv"
)

func DeHamming7(file []byte, fixErrors bool) (ret []byte) {
	var encoded1stByte, encoded2ndByte, bitsToSpare, decoded1stByte, decoded2ndByte, decodedByte byte
	bitsToSpare = 0
	two55 := exp(8) - 1 // 255
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
		decoded1stByte = decode7(encoded1stByte, fixErrors) << 4
		//Save bits that does not belong to the PracticoDeMaquina block
		bitsToSpare = file[i] & (exp(j+1) - 1)
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
			//Select second PracticoDeMaquina block
			encoded2ndByte = file[i+1] & (two55 << (j + 1))
			//Move the slice of block to its position
			encoded2ndByte = encoded2ndByte >> (j)
			//Append bits to spare and the bits that belongs to the second PracticoDeMaquina block
			encoded2ndByte = bitsToSpare | encoded2ndByte
			encoded2ndByte >>= 1
			//Save bits that does not belong to the PracticoDeMaquina block for the next iteration
			bitsToSpare = file[i+1] & (exp(j+1) - 1)
			//Append 2nd decoded half to decodedByte
			decoded2ndByte = decode7(encoded2ndByte, fixErrors)
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

func decode7(bait byte, fixErrors bool) (s byte) {
	c1 := (bait & uint8(64)) >> 6
	c2 := (bait & uint8(32)) >> 5
	d1 := (bait & uint8(16)) >> 4
	c3 := (bait & uint8(8)) >> 3
	d2 := (bait & uint8(4)) >> 2
	d3 := (bait & uint8(2)) >> 1
	d4 := (bait & uint8(1)) >> 0
	if fixErrors {
		//Calculate sindrome using xor
		var s1, s2, s3 byte

		s1 = (c1 ^ d1 ^ d2 ^ d4) << 0
		s2 = (c2 ^ d1 ^ d3 ^ d4) << 1
		s3 = (c3 ^ d2 ^ d3 ^ d4) << 2

		s = s1 | s2 | s3

		if s != 0 {
			bait = correct(bait, s)
		}
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

//correct Corrects the bit containing the error.
//
// bait: The PracticoDeMaquina block for PracticoDeMaquina 7. Extra 0 is assume to be at the right.
//
// syndrome : Position (left to right) where the mistake is.
//
// returns the block without the error.
func correct(bait byte, syndrome byte) (corrected byte) {
	//Get the wrong bit
	mistake := bait & exp(7-syndrome)
	//If it is 0, the only way to know its position is using the same power of 2
	if mistake == 0 {
		mistake = exp(7 - syndrome)
	} else { //If the bit is 1 it has to be 0
		mistake = 0
	}
	//wom comes from Without Mistake, which is bait with 0 in the position of the mistake
	wom := bait & (255 - exp(7-syndrome))

	corrected = wom | mistake

	return corrected
}

//Check errors, invoke to the function decode for decoding all the input file and finally compress the result when it's necessary (32 and  1024 bits)
func CallDecode(size int, input []byte, fixErrors bool) []byte {
	lenInput, _ := strconv.ParseInt(string(input[len(input)-10:]), 10, 64)
	if lenInput == 0 {
		return []byte{}
	}
	input = input[:len(input)-10]
	var decodedFile []byte
	_, _, controlBitsQuantity := initialCase(size)
	blockSize := size / 8
	il := 0
	sl := blockSize
	//For every block of 32, 1024 or 32768 bits check errors and then decode it
	for i := 0; i < len(input); i += blockSize {
		if fixErrors {
			checkError(size, input[il:sl], controlBitsQuantity)
		}
		aux := decode(size, input[il:sl], controlBitsQuantity)
		for j := 0; j < len(aux); j++ {
			decodedFile = append(decodedFile, aux[j])
		}
		il += blockSize
		sl += blockSize
	}
	//Finally compress
	var ret []byte
	switch size {
	case 32:
		ret = compress32(decodedFile)
	case 1024:
		ret = compress1024(decodedFile)
	case 32768:
		ret = decodedFile
	}
	return ret[:lenInput]
}

//Take the data bits from the hamming block of 32,1024 and 32768 bits
func decode(size int, input []byte, controlBitsQuantity int) []byte {
	decoded := make([]byte, int(math.Ceil((float64(size)-float64(controlBitsQuantity))/8)))
	//Take control of the position of the bit who is being seeing by the function
	decodedPosition := 0
	//Take control of the byte of the PracticoDeMaquina block who is being seeing by the function
	decodedNumberOfByte := 0
	for i := 0; i < controlBitsQuantity-1; i++ {
		//Position of the lower control bit (inferior limit)
		il := expInt(i) - 1
		//Position of the higher control bit (superior limit)
		sl := expInt(i+1) - 1
		//Take every bit between the two control bits calculated recently
		for j := il + 1; j < sl; j++ {
			position, numberOfByte := byteNumberDeHamming(j)
			//Take the data bit in the current position
			dataBit := takeBit(input[numberOfByte], 7-position, 7-decodedPosition)
			//Insert it in the return array
			decoded[decodedNumberOfByte] = decoded[decodedNumberOfByte] | dataBit
			decodedPosition++
			if decodedPosition > 7 {
				decodedPosition = 0
				decodedNumberOfByte++
			}
		}
	}
	//Add a zero to the decoded file to facilitate the posterior compression. In 32 bits and 32768 bits hamming it's not necessary
	if size == 1024 {
		decoded = append(decoded, 0)
	}
	return decoded
}

//Return the position mod 8 and the number of byte where belongs the position
func byteNumberDeHamming(position int) (int, int) {
	place := position % 8
	byteNumber := position / 8
	return place, byteNumber
}

func compress32(input []byte) []byte {
	have := make([]int, len(input)/4)
	need := make([]int, len(input)/4)
	input = append(input, 0)
	var compressed []byte //:= make([]byte, int(math.Ceil(float64(len(input))*3.25/float64(4))))
	for i := 0; i < len(input)/4; i++ {
		have[i] = 26
		need[i] = 6
	}
	position := 2
	for i := 0; i < (len(input)-1)/4-1; i++ {
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
		} else if need[i] > have[i+1] && i != (len(input)-1)/4-2 {
			aux1 := takeBitsDeHamming(have[i+1], input[(i+1)*4:(i+1)*4+5], position)
			aux2 := takeBitsDeHamming(need[i]-have[i+1], input[(i+2)*4:(i+2)*4+5], (position+have[i+1])%8)
			conflictByte1 := int(have[i] / 8)
			for j := 0; j < conflictByte1; j++ {
				compressed = append(compressed, input[i*4+j])
			}
			compressed = append(compressed, input[i*4+conflictByte1]|aux1[0])
			for j := conflictByte1 + 1; j < len(aux1)-1; j++ {
				compressed = append(compressed, aux1[j])
			}
			conflictByte2 := int(have[i+1] / 8)
			compressed = append(compressed, aux1[conflictByte2]|aux2[0])
			for j := 1; j < len(aux2); j++ {
				compressed = append(compressed, aux2[j])
			}
			aux := ajustBytes(32, input[(i+2)*4:(i+2)*4+5], need[i]-have[i+1])
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
				compressed = append(compressed, input[i*4+j])
			}
			compressed = append(compressed, input[i*4+conflictByte]|aux[0])
			for j := 1; j < len(aux); j++ {
				compressed /*[i*4+j]*/ = append(compressed, aux[j])
			}
			need[i+1] += need[i]
			have[i+1] -= need[i]
			aux = ajustBytes(32, input[(i+1)*4:(i+1)*4+4+1], need[i])
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
		if input[len(input)-5+i] != 0 {
			compressed = append(compressed, input[len(input)-5+i])
		} else {
			compressed = append(compressed, input[len(input)-5+i])
		}
	}
	return compressed
}

func compress1024(input []byte) []byte {
	bytesBlock := 128
	have := make([]int, (len(input)+1)/bytesBlock)
	need := make([]int, (len(input)+1)/bytesBlock)
	input = append(input, 0)
	var compressed []byte
	for i := 0; i < len(have); i++ {
		have[i] = 1013
		need[i] = 11
	}
	position := 5
	for i := 0; i < (len(input)-1)/bytesBlock-1; i++ {
		if need[i]%8 == 0 && need[i] <= 1013 {
			position = 5
			for j := 0; j < have[i]/8; j++ {
				compressed = append(compressed, input[i*bytesBlock+j])
			}
			for j := 0; j < need[i]/8; j++ {
				compressed = append(compressed, input[(i+1)*bytesBlock+j])
			}
			need[i+1] += need[i]
			have[i+1] -= need[i]
			index := 0
			for j := need[i] / 8; j < bytesBlock; j++ {
				input[(i+1)*bytesBlock+index] = input[(i+1)*bytesBlock+j]
				index++
			}
			for j := index; j < bytesBlock; j++ {
				input[(i+1)*bytesBlock+j] = byte(0)
			}
		} else if need[i] > have[i+1] /*i != (len(input)-1)/bytesBlock-2*/ {
			aux1 := takeBitsDeHamming(have[i+1], input[(i+1)*bytesBlock:(i+1)*bytesBlock+bytesBlock+1], position)
			aux2 := takeBitsDeHamming(need[i]-have[i+1], input[(i+2)*bytesBlock:(i+2)*bytesBlock+bytesBlock+1], (position+have[i+1])%8)
			var conflictByte2 int
			if need[i]%8 != 0 {
				conflictByte1 := int(have[i] / 8)
				for j := 0; j < conflictByte1; j++ {
					compressed = append(compressed, input[i*bytesBlock+j])
				}
				compressed = append(compressed, input[i*bytesBlock+conflictByte1]|aux1[0])
				for j := 1; j < len(aux1)-1; j++ {
					compressed = append(compressed, aux1[j])
				}
				conflictByte2 = int(have[i+1] / 8)
				if position >= 4 && position <= 7 {
					conflictByte2++
				}
				if (position+have[i+1])%8 != 0 {
					compressed = append(compressed, aux1[conflictByte2]|aux2[0])
					for j := 1; j < len(aux2); j++ {
						compressed = append(compressed, aux2[j])
					}
				} else {
					compressed = append(compressed, aux1[len(aux1)-1])
					for j := 0; j < len(aux2); j++ {
						compressed = append(compressed, aux2[j])
					}
				}
			} else {
				compressed = append(compressed, input[i*bytesBlock])
				for j := 0; j < bytesBlock-2; j++ {
					compressed = append(compressed, aux1[j])
				}
				compressed = append(compressed, aux1[126]|aux2[0])
			}
			aux := ajustBytes(1024, input[(i+2)*bytesBlock:(i+2)*bytesBlock+bytesBlock+1], need[i]-have[i+1])
			for j := 0; j < len(aux); j++ {
				input[(i+2)*bytesBlock+j] = aux[j]
			}
			have[i+2] -= need[i] - have[i+1]
			need[i+2] += need[i] - have[i+1]
			i++
			position = (position + 2) % 8

		} else {
			aux := takeBitsDeHamming(need[i], input[(i+1)*bytesBlock:(i+1)*bytesBlock+bytesBlock+1], position)
			conflictByte := int(have[i] / 8)
			for j := 0; j < conflictByte; j++ {
				compressed = append(compressed, input[i*bytesBlock+j])
			}
			compressed = append(compressed, input[i*bytesBlock+conflictByte]|aux[0])
			for j := 1; j < len(aux); j++ {
				compressed = append(compressed, aux[j])
			}
			need[i+1] += need[i]
			have[i+1] -= need[i]
			aux = ajustBytes(1024, input[(i+1)*bytesBlock:(i+1)*bytesBlock+bytesBlock+1], need[i])
			for j := 0; j < len(aux); j++ {
				input[(i+1)*bytesBlock+j] = aux[j]
			}
			if need[i] == 1013 {
				i++
				position = (position + 2) % 8
			} else {
				position = (position + 5) % 8
			}
		}
	}
	for i := 0; i < bytesBlock; i++ {
		if len(input) != bytesBlock+1 {
			if input[len(input)-(bytesBlock+1)+i] != 0 {
				compressed = append(compressed, input[len(input)-(bytesBlock+1)+i])
			}
		} else {
			compressed = append(compressed, input[len(input)-(bytesBlock+1)+i])
		}
	}
	return compressed
}

//takeBitsDeHamming
//bits is the amount of bits you need.
//input is the original byte slice.
//initialPosition are the left shift to apply
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

//begin: the number of byte where begin the information bits
func ajustBytes(size int, input []byte, begin int) []byte {
	position := begin % 8
	bytesQuantity := size / 8
	numberOfByte := int(begin / 8)
	ret := make([]byte, bytesQuantity)
	ret[0] = ((doMask(8-position) >> byte(position)) & input[numberOfByte]) << byte(position)
	aux := takeBitsDeHamming((len(input)-numberOfByte-2)*8, input[numberOfByte+1:], 8-position)
	ret[0] |= aux[0]
	for i := 1; i < len(aux); i++ {
		ret[i] = aux[i]
	}
	return ret
}

func checkError(size int, input []byte, controlBitsQuantity int) {
	syndrome := make([]byte, controlBitsQuantity)
	//Control bits calculus process
	for i := 0; i < controlBitsQuantity-1; i++ {
		parity := byte(0)
		for j := expInt(i) - 1; j < size; j += expInt(i + 1) {
			for k := 0; k < expInt(i); k++ {
				parity ^= takeBit(input[int((j+int(k))/8)], 7-(int((j+k)%8)), 0)
			}
		}
		syndrome[i] = parity
	}
	syndrome[controlBitsQuantity-1] = takeBit(input[int((size-1)/8)], 0, 0)
	correct := true
	for i := 0; i < len(syndrome); i++ {
		if syndrome[i] == 1 {
			correct = false
			break
		}
	}
	if !correct {
		errorPosition := 0
		for i := len(syndrome) - 1; i >= 0; i-- {
			errorPosition += expInt(i) * int(syndrome[i])
		}
		numberOfByte := byteNumber(errorPosition-1, size/4)
		position := (errorPosition - 1) % 8
		mistake := input[numberOfByte] & exp(byte(7-position))
		if mistake == 0 {
			mistake = exp(byte(7 - position))
		} else {
			mistake = 0
		}
		//wom comes from Without Mistake, which is bait with 0 in the position of the mistake
		wom := input[numberOfByte] & (255 - exp(byte(7-position)))
		input[numberOfByte] = wom | mistake
	}
}

/*func howManyZeros(input []byte) int {
	countZeros := 0
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == 0 {
			countZeros++
		}
	}
	return countZeros
}*/

func exp(exponent byte) (ret byte) {
	switch exponent {
	case 0:
		return 1
	case 1:
		return 2
	case 2:
		return 4
	case 3:
		return 8
	case 4:
		return 16
	case 5:
		return 32
	case 6:
		return 64
	case 7:
		return 128
	default:
		return 0
	}
}

func byteNumber(bitPosition int, bytesQuantity int) int {
	il := 0
	sl := 7
	for i := 0; i < bytesQuantity; i++ {
		if bitPosition >= il && bitPosition <= sl {
			return i
		} else {
			il += 8
			sl += 8
		}
	}
	return 0
}
