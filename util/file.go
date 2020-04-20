package util

import (
	"io/ioutil"
	"log"
	"os"
)

// CreateTmp ...
func CreateTmp(folder, file string) (string, *os.File) {

	tmpFolder, err := ioutil.TempDir("", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	f, err := ioutil.TempFile(tmpFolder, file)
	if err != nil {
		log.Fatal(err)
	}

	return tmpFolder, f

}
