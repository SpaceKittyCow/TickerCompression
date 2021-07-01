package Ticker

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//Compress compresses JSON into a CSV like format for Polygon.io data. It also does subtraction on all the data it can to shrink the data even more. TODO: This should be converted to using concurrancy for long calls.
func Compress(data string) (string, error) {
	var (
		ticker          = Ticker{}
		outliner        Ticker
		resultsToEncode []Result
		emptyResult     []Result
	)

	err := json.Unmarshal([]byte(data), &ticker)
	if err != nil {
		return "", err
	}

	resultsToEncode = ticker.Results
	log.Printf("start encoding")
	encodedResults, err := compressResults(resultsToEncode)
	if err != nil {
		return "", err
	}

	outliner = ticker
	outliner.Results = emptyResult

	outline, err := json.Marshal(outliner)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s \n%s", outline, encodedResults), nil
}

func compressResults(resultsToEncode []Result) (string, error) {
	var (
		masterList string
		lastResult = Result{}
	)

	for i := 0; i < len(resultsToEncode); i++ {
		var compress = Result{}

		if i == 0 {
			lastResult = resultsToEncode[i]
			masterList = fmt.Sprintf("%s%s", masterList,
				writeCompressedLine(lastResult))
			continue
		}

		log.Printf("at %d", i)
		compress = resultsToEncode[i]
		compress.SIP = resultsToEncode[i].SIP - lastResult.SIP
		compress.Participant = resultsToEncode[i].Participant - lastResult.Participant
		compress.Sequence = resultsToEncode[i].Sequence - lastResult.Sequence
		compress.Exchange = resultsToEncode[i].Exchange - lastResult.Exchange
		compress.Size = resultsToEncode[i].Size - lastResult.Size
		//Floating point can cause inaccuries
		//compress.Price = resultsToEncode[i].Price - lastResult.Price

		masterList = fmt.Sprintf("%s%s", masterList,
			writeCompressedLine(compress))
		lastResult = resultsToEncode[i]

	}

	return masterList, nil
}

func writeCompressedLine(compressed Result) string {
	// SIP t Participant y Sequence q ID i Exchange x Size s Conditions c Price p Tape z
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
		strconv.FormatInt(compressed.SIP, 10),
		strconv.FormatInt(compressed.Participant, 10),
		strconv.Itoa(compressed.Sequence),
		compressed.ID,
		strconv.Itoa(compressed.Exchange),
		strconv.Itoa(compressed.Size),
		fmt.Sprintf("%d", compressed.Conditions),
		strconv.FormatFloat(compressed.Price, 'f', -1, 64),
		strconv.Itoa(compressed.Tape))

}

//Decompress takes the Data from Compress and is able to reproduce the JSON that was passed in.
func Decompress(compressedData string) (string, error) {
	var (
		lastResult         = Result{}
		decompressedTicker = Ticker{}
	)

	seperatedData := strings.Split(compressedData, "\n")
	err := json.Unmarshal([]byte(seperatedData[0]), &decompressedTicker)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(seperatedData)-1; i++ {

		log.Println("%d", i)
		decompressed, err := readCompressedLine(seperatedData[i])
		if err != nil {
			return "", err
		}
		decompressed.SIP = decompressed.SIP + lastResult.SIP
		decompressed.Participant = decompressed.Participant + lastResult.Participant
		decompressed.Sequence = decompressed.Sequence + lastResult.Sequence
		decompressed.Exchange = decompressed.Exchange + lastResult.Exchange
		decompressed.Size = decompressed.Size + lastResult.Size
		//Floating point can cause inaccuries
		// decompressed.Price = decompressed.Price + lastResult.Price

		lastResult = decompressed

		decompressedTicker.Results = append(decompressedTicker.Results, decompressed)
	}
	final, err := json.Marshal(decompressedTicker)
	if err != nil {
		return "", err
	}
	return string(final), nil
}

func readCompressedLine(decompress string) (Result, error) {
	var (
		decompressed = Result{}
		err          error
	)
	// SIP t Participant y Sequence q ID i Exchange x Size s Conditions c Price p Tape z
	seperatedFields := strings.Split(decompress, ",")
	log.Printf(seperatedFields[4])
	decompressed.SIP, err = strconv.ParseInt(seperatedFields[0], 10, 64)
	if err != nil {
		return Result{}, err
	}
	decompressed.Participant, err = strconv.ParseInt(seperatedFields[1], 10, 64)
	if err != nil {
		return Result{}, err
	}
	decompressed.Sequence, err = strconv.Atoi(seperatedFields[2])
	if err != nil {
		return Result{}, err
	}
	decompressed.ID = seperatedFields[3]
	decompressed.Exchange, err = strconv.Atoi(seperatedFields[4])
	if err != nil {
		return Result{}, err
	}

	decompressed.Size, err = strconv.Atoi(seperatedFields[5])
	if err != nil {
		return Result{}, err
	}
	var inter []int
	str := strings.Replace(seperatedFields[6], " ", ",", -1)
	err = json.Unmarshal([]byte(str), &inter)

	if err != nil {
		log.Fatal(err)
	}
	decompressed.Conditions = inter
	decompressed.Price, err = strconv.ParseFloat(seperatedFields[7], 64)
	if err != nil {
		return Result{}, err
	}

	decompressed.Tape, err = strconv.Atoi(seperatedFields[8])
	if err != nil {
		return Result{}, err
	}
	return decompressed, nil
}
