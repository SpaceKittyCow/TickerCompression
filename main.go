package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"polygon.io/TickerCompression"
)

func main() {
	jsonFileLocation, compressedFileLocation, saveFileLocation := constructFlags()
	err := ValidateFiles(jsonFileLocation, compressedFileLocation, saveFileLocation)
	if err != nil {
		log.Fatal(err)
	}

	var saveData string
	if jsonFileLocation != "" {
		jsonFile, err := ImportFile(jsonFileLocation)
		if err != nil {
			log.Fatal(err)
		}
		saveData, err = TickerCompression.Compress(jsonFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if compressedFileLocation != "" {
		compressedFile, err := ImportFile(compressedFileLocation)
		if err != nil {
			log.Fatal(err)
		}
		saveData, err = TickerCompression.Decompress(compressedFile)
		if err != nil {
			log.Fatal(err)
		}

	}

	err = SaveFile(saveData, saveFileLocation)

	if err != nil {
		log.Fatal(err)
	}

}

func constructFlags() (string, string, string) {
	jsonFileLocation := flag.String("c", "", "Specify a ticker json file location")
	compressedFileLocation := flag.String("d", "", "Specify a compressed file location")
	saveFileLocation := flag.String("s", "", "Specify a final save point for file")
	flag.Parse()

	return *jsonFileLocation, *compressedFileLocation, *saveFileLocation
}

//Validates files to both see if the correct file locations are given and to see if the files esist, or in the case of a save file, doesn't exist.
func ValidateFiles(jsonFileLocation, compressedFileLocation, saveFileLocation string) error {

	if (jsonFileLocation == "" && compressedFileLocation == "") || (jsonFileLocation != "" && compressedFileLocation != "") {
		return fmt.Errorf("You must specify either a ticker json file location or a compressed file location")
	}

	if jsonFileLocation != "" {
		err := checkFileLocation(jsonFileLocation)
		if err != nil {
			return err
		}
	}

	if compressedFileLocation != "" {
		err := checkFileLocation(compressedFileLocation)
		if err != nil {
			return err
		}
	}

	err := checkFileLocation(saveFileLocation)
	if err == nil || !os.IsNotExist(err) {
		return fmt.Errorf("Save File Location already exists")
	} else if err != nil {
		return err
	}

	return nil

}

func checkFileLocation(fileLocation string) error {
	_, err := os.Stat(fileLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("A File Location does not exist. Double check file %s ", fileLocation)
		}
		return err
	}
	return nil
}

//ImportFile just imports the file into a string.
func ImportFile(fileLocation string) (string, error) {
	fileContents, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return "", err
	}

	text := string(fileContents)
	return text, nil
}

//Save File just saves the returned string to a file
func SaveFile(saveData, saveFileLocation string) error {

	file, err := os.Create(saveFileLocation)
	if err != nil {
		return err
	}
	defer file.Close()

	err = ioutil.WriteFile(saveFileLocation, []byte(saveData), 0777)
	if err != nil {
		return err
	}

	return nil
}
