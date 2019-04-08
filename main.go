package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"
	//"www.github.com/willf/bitset.git"
)

func main() {
	start := time.Now()
	var fileName string
	var encodedBody []byte

	fmt.Println("Ingrese el nombre del archivo")
	_, _ = fmt.Scanf("%s", &fileName)

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
	log.Printf("Hamming7 took %s", elapsed)
}

// Receives a byte slice, returns it encoded
func hamming7(file []byte) []byte {
	//Mask that shows first bits
	maskFirst := 240
	//Mask that shows last bits
	maskLast := 15
	//Encoded byte slice
	var ret []byte
	for i := 0; i < len(file); i++ {
		//Add mask
		maskedFirst := file[i] & uint8(maskFirst)
		//Move to last positions
		maskedFirst /= 16
		//Add mask
		maskedLast := file[i] & uint8(maskLast)

		//Convert int to byte slices
		slicedFirst := parceByte(int8(maskedFirst))
		slicedLast := parceByte(int8(maskedLast))

		//Encode with hamming 7 both parts
		slicedFirst, slicedLast = encode(slicedFirst, slicedLast)

		ret = append(ret, byte(deParceByte(slicedFirst)))
		ret = append(ret, byte(deParceByte(slicedLast)))

	}
	return ret
}

func parceByte(maskedByte int8) [7]int {
	//Convert maskedByte to a string with its binary representation
	stringByte := FormatInt(maskedByte, 2)
	//Convert stringByte to a slice
	sliceSB := []byte(stringByte)
	var ret [7]int
	sliceLen := len(sliceSB)
	for i := sliceLen - 1; i >= 0; i-- {
		//Convert strings to int in all positions
		ret[7-sliceLen+i], _ = strconv.Atoi(string(sliceSB[i]))
	}
	return ret
}

func deParceByte(encoded [7]int) int8 {
	encodedString := ""
	for i := 0; i < 7; i++ {
		encodedString = strconv.Itoa(encoded[i])
	}
	ret, _ := ParseInt(encodedString, 2, 8)
	return ret
}

func encode(target1 [7]int, target2 [7]int) ([7]int, [7]int) {
	target1[2] = target1[3]
	target1[0] = target1[2] ^ target1[4] ^ target1[6]
	target1[1] = target1[2] ^ target1[5] ^ target1[6]
	target1[3] = target1[4] ^ target1[5] ^ target1[6]

	target2[2] = target2[3]
	target2[0] = target2[2] ^ target2[4] ^ target2[6]
	target2[1] = target2[2] ^ target2[5] ^ target2[6]
	target2[3] = target2[4] ^ target2[5] ^ target2[6]

	return target1, target2
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
