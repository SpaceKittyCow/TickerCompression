package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"polygon.io/Ticker"
	"time"
)

type configuration struct {
	ApiKey                 string
	day                    int
	month                  int
	year                   int
	Date                   *time.Time
	CompressedFileLocation string
	OrignalFileLocation    string
	Stock                  string
}

func main() {
	var (
		orignalData, compressedData, decompressedData string
	)
	config, err := hydrateConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	orignalData, err = Ticker.GetDaysData(config.ApiKey, config.Stock, config.Date)
	if err != nil {
		log.Fatal(err)
	}

	_ = SaveFile(orignalData, config.OrignalFileLocation)
	compressedData, err = Ticker.Compress(orignalData)
	if err != nil {
		log.Fatal(err)
	}

	decompressedData, err = Ticker.Decompress(compressedData)
	if err != nil {
		log.Fatal(err)
	}

	if decompressedData != orignalData {
		_ = SaveFile(compressedData, config.CompressedFileLocation)

		log.Fatal(fmt.Errorf("Files do not match"))
	}

	log.Printf("Files Match")
}

func hydrateConfiguration() (configuration, error) {

	config := createConfiguration()
	if config.ApiKey == "" {
		return configuration{}, fmt.Errorf("API Key needed")
	}

	err := validateFiles(config.CompressedFileLocation, config.OrignalFileLocation)
	if err != nil {
		return configuration{}, err
	}

	config.Date, err = validateAndFormatDate(config.day, config.month, config.year)
	if err != nil {
		return configuration{}, err
	}

	return config, nil
}

func createConfiguration() configuration {

	apiKey := flag.String("a", "", "Polygon.io API Key")
	day := flag.Int("d", 0, "Specify a day of the month: 1, 2, 3 ...")
	month := flag.Int("m", 0, "Specify a month number: 1 for January, 2 for Feburary ...")
	year := flag.Int("y", 0, "Specify a year: 2021, 2020 ...")
	compressedFileLocation := flag.String("c", "./compressedfile", "Specify a compressed file location")
	orignalFileLocation := flag.String("s", "./orignalfile.json", "Specify a final save point for file")
	flag.Parse()

	return configuration{
		Stock:                  "AAPL",
		ApiKey:                 *apiKey,
		CompressedFileLocation: *compressedFileLocation,
		OrignalFileLocation:    *orignalFileLocation,
		day:                    *day,
		month:                  *month,
		year:                   *year,
	}

}

func validateAndFormatDate(day int, month int, year int) (*time.Time, error) {

	if day == 0 || month == 0 || year == 0 {
		return nil, fmt.Errorf("Please supply day, month and year")
	}

	if day > 31 {
		return nil, fmt.Errorf("Day is not a real day")
	}

	//Year Apple stock went live
	if year < 1980 {
		return nil, fmt.Errorf("Year is not possible for Stock")
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if date.After(time.Now()) {
		return nil, fmt.Errorf("Period of time has not occured yet")
	}

	return &date, nil

}
func validateFiles(compressedFileLocation, orignalFileLocation string) error {

	err := checkFileLocation(compressedFileLocation)
	if err != nil {
		return err
	}

	err = checkFileLocation(orignalFileLocation)
	if err != nil {
		return err
	}

	return nil

}

func checkFileLocation(fileLocation string) error {
	_, err := os.Stat(fileLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return fmt.Errorf("A File Location exists already. Double check file %s ", fileLocation)
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
