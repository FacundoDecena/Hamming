package HammingCodification

import (
	"math"
	"strconv"
)

// Receives a byte slice, returns it encoded
func Hamming7(file []byte) []byte {
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
	//Applies the PracticoDeMaquina encode to each byte of the file
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

func Hamming(size int, file []byte) []byte {
	var ret []byte
	switch size {
	case 32:
		x := callTakeBits(26, file)
		ret = callEncode(size, x, len(file))
	case 1024:
		x := callTakeBits(1013, file)
		ret = callEncode(size, x, len(file))
	case 32768:
		x := convertTo32752(file)
		ret = callEncode(size, x, len(file))
	}
	return ret
}

func convertTo32752(input []byte) [][]byte {
	ret := make([][]byte, int(math.Ceil(float64(len(input))/float64(4094))))
	il := 0
	var sl int
	if len(input) < 4094 {
		sl = len(input)
	} else {
		sl = 4094
	}
	breakLastPosition := false
	var i int
	for i = 0; i < len(ret); i++ {
		if sl-il == 4094 {
			ret[i] = input[il:sl]
			il += 4094
			if len(input)-(i+1)*4094 < 4094 {
				sl += len(input) - (i+1)*4094
			} else {
				sl += 4094
			}
		} else {
			breakLastPosition = true
			break
		}
	}
	if breakLastPosition {
		var lastPosition []byte
		for j := il; j < sl; j++ {
			lastPosition = append(lastPosition, input[j])
		}
		for j := 0; j < 4094-sl%4094; j++ {
			lastPosition = append(lastPosition, 0)
		}
		ret[i] = lastPosition
	}
	return ret
}

func callEncode(size int, inputFile [][]byte, lenFile int) (outPut []byte) {
	position, numberOfByte, controlBitsQuantity := initialCase(size)
	var aux [][]byte
	//matrix := getGeneratingMatrix(size)
	for i := 0; i < len(inputFile); i++ {
		aux = append(aux, encode(size, inputFile[i], position, numberOfByte, controlBitsQuantity /*,matrix*/))
		for j := 0; j < len(aux[i]); j++ {
			outPut = append(outPut, aux[i][j])
		}
	}
	s := []byte(strconv.FormatInt(int64(lenFile), 10))
	lenInput := []byte(s)
	for i := len(lenInput); i < 10; i = len(lenInput) {
		lenInput = append([]byte{48}, lenInput...)
	}
	outPut = append(outPut, lenInput...)
	return outPut
}

//Size should be: 8 for hamming7, 32 for Hamming 32, 1024 for Hamming 1024 and 32768 for Hamming 32768
func encode(size int, input []byte, position int, numberOfByte int, controlBitsQuantity int) []byte {
	encoded := make([]byte, int(size/8))
	//Data bits accommodate process
	for i := controlBitsQuantity - 1; i > 0; i-- {
		sl := expInt(i) - 1
		il := expInt(i-1) - 1
		for j := sl - 1; j > il; j-- {
			dataBit := takeBit(input[numberOfByte], position, 7-int(j%8))
			x := int(j / 8)
			encoded[x] = encoded[x] | dataBit
			position++
			if position > 7 {
				numberOfByte--
				position = 0
			}
		}
	}
	//Control bits calculus process
	for i := 0; i < controlBitsQuantity-1; i++ {
		parity := byte(0)
		for j := expInt(i) - 1; j < size; j += expInt(i + 1) {
			for k := 0; k < expInt(i); k++ {
				parity ^= takeBit(encoded[int((j+k)/8)], 7-((j+k)%8), 0)
			}
		}
		x := int(int(expInt(i)-1) / 8)
		encoded[x] = encoded[x] | takeBit(parity, 0, 7-(expInt(i)-1)%8)
	}
	return encoded
}

//Matrix implementation for encode
/*
func encode(size int, input []byte, position int, numberOfByte int, controlBitsQuantity int, matrix[] byte) []byte {
	encoded := make([]byte, int(size/8))
	controlBits := getControlBits(input,matrix,controlBitsQuantity,size)
	//Bits accommodate process
	for i := controlBitsQuantity - 1; i > 0; i-- {
		sl := expInt(i) - 1
		il := expInt(i-1) - 1
		m := int(sl/8)
		encoded[m] = encoded[m] | controlBits[i] << byte(7-(sl%8))
		for j := sl - 1; j > il; j-- {
			dataBit := takeBit(input[numberOfByte], position, 7-int(j%8))
			x := int(j/8)
			encoded[x] = encoded[x] | dataBit
			position++
			if position > 7 {
				numberOfByte--
				position = 0
			}
		}
	}
	encoded[0] = encoded[0] | (controlBits[0]<<7)
	return encoded
}

func getControlBits(vector []byte,matrix []byte,controlBitsQuantity int,size int) []byte{
	controlBits := make([]byte,controlBitsQuantity)
	for i:=0;i<controlBitsQuantity-1;i++{
		parity := byte(0)
		for j:=0;j<size-controlBitsQuantity;j++{
			x := int(int(j)/8)
			parity ^= (vector[x] & byte(expInt(7-int(j%8))) )>> byte(7-(j%8)) & matrix[j*(controlBitsQuantity-1)+i]
		}
		controlBits[i]=parity
	}
	return controlBits
}


func getGeneratingMatrix(size int) []byte{
	_,_,m :=initialCase(size)
	matrix := make([]byte,0)
	var column int
	controlBit:=0
	for row:=0;row<size-1;row++{
		if row != expInt(controlBit)-1{
			for column=0;column<m-1;column++ {
				matrix = append(matrix, b2b(getBoolean(column, row)))
			}
			column=0
		}else{
			controlBit++
		}
	}
	return matrix
}

func getBoolean(i int,j int) bool{
	switch i{
	case 0:return j%2==0
	case 1:return j%4>=1 && j%4<=2
	case 2:return j%8>=3 && j%8<=6
	case 3:return j%16>=7 && j%16<=14
	case 4:return j%32>=15 && j%32<=30
	case 5:return j%64>=31 && j%64<=62
	case 6:return j%128>=63 && j%128<=126
	case 7:return j%256>=127 && j%256<=254
	case 8:return j%512>=255 && j%512<=510
	case 9:return j%1024>=511 && j%1024<=1022
	case 10:return j%2048>=1023 && j%2048<=2046
	case 11:return j%4096>=2047 && j%4096<=4094
	case 12:return j%8192>=4095 && j%8192<=8190
	case 13:return j%16384>=8191 && j%16384<=16382
	case 14:return j%32768>=16383 && j%32768<=32766
	default: return false
	}
}

func b2b(b bool) byte {
	if b{
		return 1
	}
	return 0
}

*/

//Apply a mask to a source byte to get the bit in the initial position and shifter it to the final position.
func takeBit(source byte, initialPosition int, finalPosition int) byte {
	result := source & byte(expInt(initialPosition))
	shift := finalPosition - initialPosition
	if shift >= 0 {
		return result << uint(shift)
	} else {
		return result >> uint(shift*-1)
	}
}

//Implementation for exponential base 2. It's faster than calculate exponential with iterations
func expInt(exponent int) int {
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
	case 8:
		return 256
	case 9:
		return 512
	case 10:
		return 1024
	case 11:
		return 2048
	case 12:
		return 4096
	case 13:
		return 8192
	case 14:
		return 16384
	case 15:
		return 32768
	default:
		return -1
	}
}

func initialCase(size int) (position int, numberOfByte int, controlBitsQuantity int) {
	//Set the initial position where is the first information bit in the array passed by parameter depending of what PracticoDeMaquina will be apply
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
		numberOfByte = 4093
		controlBitsQuantity = 16
	}
	return position, numberOfByte, controlBitsQuantity
}
