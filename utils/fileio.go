package utils

import (
	"os"
	"log"
	"io/ioutil"
)

func readFile(file string) ([]byte, error) {
	reader, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	return data, err
}
