package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
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

func actualTime() (time.Time, error) {
	cmd := exec.Command("curl", "http://worldtimeapi.org/api/timezone/America/Argentina/Cordoba")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return time.Time{}, errors.New("No se pudo establecer conexion con el servidor. Revise su conexion.")
	} else {
		var a Answer
		err = json.Unmarshal(out.Bytes(), &a)
		if err != nil {
			return time.Time{}, errors.New("Existe un problema con el servidor. Contacte a los desarrolladores.")
		} else {
			return time.Unix(a.UnixTime, 0), nil
		}
	}
}
