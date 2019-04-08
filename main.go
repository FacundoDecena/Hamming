package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

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
	//TODO: optimize routines, learn goroutines?
	ch1 := make(chan [7]int)
	ch2 := make(chan [7]int)
	//Mask that shows first bits
	maskFirst := 240
	//Mask that shows last bits
	maskLast := 15
	//Encoded byte slice
	var ret []byte
	for i := 0; i < len(file); i++ {
		defer wg.Add(2)
		//Add mask
		maskedFirst := file[i] & uint8(maskFirst)
		//Move to last positions
		maskedFirst /= 16
		//Add mask
		maskedLast := file[i] & uint8(maskLast)

		//Convert int to byte slices

		go func() {
			stringByte := FormatInt8(int8(maskedFirst), 2)
			//Convert stringByte to a slice
			sliceSB := []byte(stringByte)
			var ret [7]int
			sliceLen := len(sliceSB)
			for i := sliceLen - 1; i >= 0; i-- {
				//Convert strings to int in all positions
				ret[7-sliceLen+i], _ = strconv.Atoi(string(sliceSB[i]))
			}
			ch1 <- ret
			defer wg.Done()
		}()
		go func() {
			stringByte := FormatInt8(int8(maskedLast), 2)
			//Convert stringByte to a slice
			sliceSB := []byte(stringByte)
			var ret [7]int
			sliceLen := len(sliceSB)
			for i := sliceLen - 1; i >= 0; i-- {
				//Convert strings to int in all positions
				ret[7-sliceLen+i], _ = strconv.Atoi(string(sliceSB[i]))
			}
			ch2 <- ret
			defer wg.Done()
		}()
		slicedFirst := <-ch1
		slicedLast := <-ch2
		wg.Wait()

		//TODO: Apply concurrency on this functions
		//Encode with hamming 7 both parts
		slicedFirst, slicedLast = encode(slicedFirst, slicedLast)

		ret = append(ret, byte(deParseByte(slicedFirst)))
		ret = append(ret, byte(deParseByte(slicedLast)))

	}
	return ret
}

func deParseByte(encoded [7]int) int8 {
	sliceLen := len(encoded)
	var encodedStringSlice []string
	for i := 0; i < sliceLen; i++ {
		encodedStringSlice = append(encodedStringSlice, strconv.Itoa(encoded[i]))
	}
	encodedString := strings.Join([]string(encodedStringSlice), "")
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
