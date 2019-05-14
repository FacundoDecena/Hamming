package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
)

func IntroduceErrors() {
	var fileName string
	var body, fileWithErrors []byte
	var err error
	r := bufio.NewReader(os.Stdin)
	clearScreen()
	fmt.Println("Ingrese el nombre del archivo a introducir errores extension:")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	//Clean buffer
	_, _ = fmt.Fscanf(r, "%s")
	body, err = loadFile(fileName)
	if err != nil {
		fmt.Println(err)
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	//Split the string between name and extension
	extension := strings.Split(fileName, ".")
	switch extension[1] {
	case "ha1":
		fileWithErrors = insertError7(body)
	case "ha2":
		fileWithErrors = insertError(body, 32)
	case "ha3":
		fileWithErrors = insertError(body, 1024)
	case "ha4":
		fileWithErrors = insertError(body, 32768)
	default:
		fmt.Println("La extension del archivo no es válida.")
		_, _ = fmt.Fscanf(r, "%s")
		return
	}
	_ = saveFile(strings.Replace(fileName, ".ha", ".he", -1), fileWithErrors)
	fmt.Println("Se han introducido errores de manera correcta.")
	_, _ = fmt.Fscanf(r, "%s")
}

func insertError7(file []byte) (ret []byte) {
	var encoded1stByte, encoded2ndByte, bitsToSpare, errored1stByte, errored2ndByte byte
	bitsToSpare = 0
	two55 := byte(exp(8)) - 1 // 255
	var j byte
	j = 0
	for i := 0; i < len(file); i += 2 {
		//Select the first 7-j bits
		encoded1stByte = file[i] & (two55 << (j + 1))
		//Move them to theirs position
		encoded1stByte >>= j
		//Join the pieces
		encoded1stByte = bitsToSpare | encoded1stByte
		//Append decoded half to decodedByte
		errored1stByte = randomErrors7(encoded1stByte)
		//Save bits that does not belong to the PracticoDeMaquina block
		bitsToSpare = file[i] & (byte(exp(j+1)) - 1)
		j++
		if j%7 == 0 && i > 0 {
			i--
			//errored1stByte = errored1stByte | bitsToSpare
		}
		if i+1 == len(file) {
			errored1stByte = errored1stByte | bitsToSpare
			ret = append(ret, errored1stByte)
		} else {
			ret = append(ret, errored1stByte)
			//Move bits to their place
			bitsToSpare = bitsToSpare << (8 - j)
			//Select second PracticoDeMaquina block
			encoded2ndByte = file[i+1] & (two55 << (j + 1))
			//Move the slice of block to its position
			encoded2ndByte = encoded2ndByte >> (j)
			//Append bits to spare and the bits that belongs to the second PracticoDeMaquina block
			encoded2ndByte = bitsToSpare | encoded2ndByte
			//Save bits that does not belong to the PracticoDeMaquina block for the next iteration
			bitsToSpare = file[i+1] & (byte(exp(j+1)) - 1)
			//Append 2nd decoded half to decodedByte
			errored2ndByte = randomErrors7(encoded2ndByte)
			ret = append(ret, errored2ndByte)
		}
		j++
		bitsToSpare = bitsToSpare << (8 - j)
		if j > 7 {
			j = 0
			bitsToSpare = 0
		}
	}
	ret = compress7(ret)
	return ret
}

func compress7(file []byte) []byte {
	entryLength := len(file)
	module := 0
	if 2*entryLength%8 != 0 {
		module = 8 - 2*entryLength%8
	}
	auxLength := 2*entryLength + module
	//Number that I use so that the size of the array is a multiple of 8, thus making compression simpler
	finalLength := int(math.Ceil(float64(entryLength) * 0.875))
	j := 0
	ret := make([]byte, auxLength)
	//Compress the array
	for i := 0; i < entryLength; i += 8 {
		sevenBlock := compressBlock(file[i : i+8])
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

func insertError(file []byte, kind int) (ret []byte) {
	var blocks [][]byte
	var erroredBlock []byte
	blocks = takeBlocks(file, kind)
	for i := 0; i < len(blocks); i++ {
		erroredBlock = randomErrors(blocks[i], kind)
		for j := 0; j < (kind / 8); j++ {
			ret = append(ret, erroredBlock[j])
		}
	}

	return ret[:len(file)]
}

func randomErrors7(bait byte) (ret byte) {
	//Error does not correspond
	if rand.Intn(2) == 0 {
		ret = bait
	} else {
		//Select random position
		position := byte(rand.Intn(6)) + 1
		//Creates a mask
		mask := exp(position)
		//Get the targeted bit
		target := bait & mask
		//Introduce the error
		if target == 0 {
			target = mask
		} else { //If the bit is 1 it has to be 0
			target = 0
		}
		////wom comes from Without Mistake, which is bait with 0 in the position of the target
		wom := bait & (255 - mask)
		ret = wom | target
	}
	return ret
}

//randomErrors receives an array with a Hamming block and returns it with an error in it
func randomErrors(input []byte, kind int) (ret []byte) {
	if rand.Intn(2) == 0 { //Error does not correspond
		ret = input
	} else {
		blockSize := kind / 8
		//Select random byte
		targetedByte := rand.Intn(blockSize - 1)
		bait := input[targetedByte]
		//Select random position
		position := byte(rand.Intn(7))
		//Creates a mask
		mask := exp(position)
		//Get the targeted bit
		target := bait & mask
		//Introduce the error
		if target == 0 {
			target = mask
		} else { //If the bit is 1 it has to be 0
			target = 0
		}
		////wom comes from Without Mistake, which is bait with 0 in the position of the target
		wom := bait & (255 - mask)
		bait = wom | target
		input[targetedByte] = bait
		ret = input
	}
	return input
}

//takeBlocks returns separated PracticoDeMaquina blocks
//
//input is the file
//
//kind is the type of PracticoDeMaquina (32, 1024, 32768)
//
//return an array with PracticoDeMaquina blocks with size kind
func takeBlocks(input []byte, kind int) (ret [][]byte) {
	var length, blockSize int
	length = len(input)
	blockSize = kind / 8
	for i := 0; i < length; i += blockSize {
		tempArray := make([]byte, blockSize)
		for j := 0; j < blockSize; j++ {
			if i+j == length {
				continue
			}
			tempArray[j] = input[i+j]
		}
		ret = append(ret, tempArray)
	}
	return ret
}
