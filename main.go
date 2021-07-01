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
	ApiKey                   string
	day                      int
	month                    int
	year                     int
	Date                     *time.Time
	ResultCount              int
	CompressedFileLocation   string
	OrignalFileLocation      string
	DecompressedFileLocation string
	JSONFileLocation         string
	Stock                    string
}

const (
	CompressedFileDefault           = "./compressedfile"
	DecompressedFileDefaultLocation = "./decompressed.json"
	OrignalFileDefault              = "./orignalfile.json"
)

func main() {
	var (
		orignalData, compressedData, decompressedData string
	)
	config, err := initConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	if config.JSONFileLocation == "" {
		log.Println("Getting Ticker Data")
		orignalData, err = Ticker.GetDaysData(config.ApiKey, config.Stock, config.Date, config.ResultCount)
		if err != nil {
			log.Fatal(err)
		}
		_ = SaveFile(orignalData, config.OrignalFileLocation)
	} else {
		log.Println("Loading JSON File")
		orignalData, err = ImportFile(config.JSONFileLocation)
		if err != nil {
			log.Fatal(err)
		}

	}

	log.Println("Compressing Data")
	compressedData, err = Ticker.Compress(orignalData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Saving Compressed Data")
	_ = SaveFile(compressedData, config.CompressedFileLocation)

	decompressedData, err = Ticker.Decompress(compressedData)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", decompressedData)
	if decompressedData != orignalData {
		_ = SaveFile(decompressedData, config.DecompressedFileLocation)

		log.Fatal(fmt.Errorf("Files do not match"))
	}

	log.Printf("Files Match")
}

func initConfiguration() (configuration, error) {

	config := createConfiguration()
	if config.ApiKey == "" {
		return configuration{}, fmt.Errorf("API Key needed")
	}

	err := validateFiles(config.CompressedFileLocation, config.OrignalFileLocation, config.JSONFileLocation)
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
	resultCount := flag.Int("r", 0, "Will go through one pass through instead of all day. Reccomended due to getting the whole day taking forever. Upto 50000")
	compressedFileLocation := flag.String("c", CompressedFileDefault, "Specify a save location of the compressed file")
	orignalFileLocation := flag.String("s", OrignalFileDefault, "Specify a save location of the orignal file")
	jSONFileLocation := flag.String("j", "", "Specify a load location of an orignal file")
	flag.Parse()

	return configuration{
		Stock:                    "AAPL",
		ApiKey:                   *apiKey,
		CompressedFileLocation:   *compressedFileLocation,
		DecompressedFileLocation: DecompressedFileDefaultLocation,
		OrignalFileLocation:      *orignalFileLocation,
		ResultCount:              *resultCount,
		JSONFileLocation:         *jSONFileLocation,
		day:                      *day,
		month:                    *month,
		year:                     *year,
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

func validateFiles(compressedFileLocation, orignalFileLocation, JSONFileLocation string) error {

	err := checkFileDoesNotExist(compressedFileLocation, CompressedFileDefault)
	if err != nil {
		return err
	}

	err = checkFileDoesNotExist(orignalFileLocation, OrignalFileDefault)
	if err != nil {
		return err
	}

	err = checkFileDoesNotExist(DecompressedFileDefaultLocation, DecompressedFileDefaultLocation)
	if err != nil {
		return err
	}

	if JSONFileLocation != "" {
		_, err = os.Stat(JSONFileLocation)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("JSON file does not exist")
			}
		}
	}

	return nil

}

func checkFileDoesNotExist(fileLocation, defaultLocation string) error {
	err := os.Remove(defaultLocation)
	if err != nil {
		fmt.Println(err)
	}

	_, err = os.Stat(fileLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return fmt.Errorf("File %s already exists. Please delete and try again", fileLocation)

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
