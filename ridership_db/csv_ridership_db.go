package ridershipDB

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type CsvRidershipDB struct {
	idIdxMap      map[string]int
	csvFile       *os.File
	csvReader     *csv.Reader
	num_intervals int
}

func (c *CsvRidershipDB) Open(filePath string) error {
	c.num_intervals = 9

	// Create a map that maps MBTA's time period ids to indexes in the slice
	c.idIdxMap = make(map[string]int)
	for i := 1; i <= c.num_intervals; i++ {
		timePeriodID := fmt.Sprintf("time_period_%02d", i)
		c.idIdxMap[timePeriodID] = i - 1
	}

	// create csv reader
	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	c.csvFile = csvFile
	c.csvReader = csv.NewReader(c.csvFile)

	return nil
}

// TODO: some code goes here
// Implement the remaining RidershipDB methods
func (c *CsvRidershipDB) GetRidership(lineId string) ([]int64, error) {
	// zero out all boarding metrics
	boardings := make([]int64, c.num_intervals)
	for i := 0; i < c.num_intervals; i++ {
		boardings[i] = 0
	}

	// skip the header
	_, err := c.csvReader.Read()
	if err != nil { // return the error
		return nil, err
	}

	for {
		value, err := c.csvReader.Read()
		if err == io.EOF { // end of file, end iteration
			break
		}
		if err != nil { // return the error
			return nil, err
		}
		if value[0] != lineId { // not the line we're looking for
			continue
		}

		passengers, err := strconv.ParseInt(value[4], 10, 64)
		if err != nil {
			return nil, err
		}

		idx := c.idIdxMap[value[2]]
		boardings[idx] += passengers
	}

	// restart reader to the top
	_, err = c.csvFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	c.csvReader = csv.NewReader(c.csvFile)
	c.csvReader.FieldsPerRecord = 5

	return boardings, nil
}

func (c *CsvRidershipDB) Close() error {
	return c.csvFile.Close()
}
