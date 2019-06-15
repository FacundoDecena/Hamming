package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Answer struct {
	UnixTime int64 `json:"unixtime"`
}

func loadFile(fileName string, dateCheck bool) ([]byte, error) {
	var err error
	var body []byte
	body, err = ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if dateCheck {
		unixTime, err := strconv.ParseInt(string(body[len(body)-10:]), 10, 64)
		if err != nil {
			return nil, err
		}
		dateEnabled := time.Unix(unixTime, 0)
		actualTime, err := actualTime()
		if err != nil {
			return nil, err
		}
		if dateEnabled.After(actualTime) {
			err = errors.New("El archivo todavia no puede ser abierto. Fecha de apertura: " + dateEnabled.String())
			return nil, err
		}
		return body[:len(body)-10], nil
	} else {
		return body, nil
	}
}

func saveFile(fileName string, body []byte) error {
	return ioutil.WriteFile(fileName, body, 0600)
}

func seeSize() {
	extensions := []string{".txt", ".ha1", ".ha2", ".ha3", ".ha4", ".huf"}
	var fileName string
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el nombre del archivo sin extension.")
	_, _ = fmt.Fscanf(r, "%s", &fileName)

	for index := 0; index < len(extensions); index++ {
		body, err := loadFile(fileName+extensions[index], false)
		if err != nil {
			fmt.Print("\n", err)
		} else {
			switch extensions[index] {
			case ".txt":
				fmt.Print("El archivo inicial tiene un tamaño de:", len(body), " Bytes ", " o ", (len(body))/1024, " KB")
			case ".ha1":
				fmt.Print("\n\n Hamming 7 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha2":
				fmt.Print("\n\n Hamming 32 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha3":
				fmt.Print("\n\n Hamming 1024 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".ha4":
				fmt.Print("\n\n Hamming 32"+
					"768 tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			case ".huf":
				fmt.Print("\n\n HuffmanCodification tiene un tamaño de: ", len(body), " Bytes ", " o ", len(body)/1024, " KB")
			}
		}
	}
	fmt.Println("\n\n Presione enter para continuar")
	_, _ = fmt.Fscanf(r, "%s", &fileName)
	_, _ = fmt.Fscanf(r, "%s", &fileName)
}

func actualTime() (time.Time, error) {
	cmd := exec.Command("curl", "http://worldtimeapi.org/api/timezone/America/Argentina/Cordoba")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return time.Time{}, err
	} else {
		var a Answer
		err = json.Unmarshal(out.Bytes(), &a)
		if err != nil {
			return time.Time{}, err
		} else {
			return time.Unix(a.UnixTime, 0), nil
		}
	}
}
