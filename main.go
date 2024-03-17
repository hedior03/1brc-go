package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

type cityAggregate struct {
	min   float64
	max   float64
	sum   float64
	count int
}

type cityResult struct {
	min     float64
	max     float64
	average float64
}

const mapAllocation = 10000
const channelBufferSize = 1000
const batchSize = 1000

func main() {
	/*
		goroutine iterating through file
			read line
			bundle lines in array
			send them to the Parsing func through a channel
		goroutine iterating through the lines sent through the channel
			parse into float
			aggregate and store in map
		consolidate result (iteration)
		output result (iteration)
	*/

	if len(os.Args) == 0 {
		panic("No input file provided.")
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		panic("File couldn't be opened")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	citiesMap := make(map[string]cityAggregate, mapAllocation)

	var wg sync.WaitGroup
	linesChannel := make(chan []string, channelBufferSize)

	wg.Add(1)
	go readLines(scanner, linesChannel, &wg)

	wg.Add(1)
	go parseLines(citiesMap, linesChannel, &wg)

	wg.Wait()

	citiesResultMap := make(map[string]cityResult, mapAllocation)
	for key := range citiesMap {
		averageTemperature := citiesMap[key].sum / float64(citiesMap[key].count)
		currentCityResult := cityResult{
			min:     approxToDecimal(citiesMap[key].min),
			max:     approxToDecimal(citiesMap[key].max),
			average: approxToDecimal(averageTemperature),
		}

		citiesResultMap[key] = currentCityResult
	}

	for key := range citiesResultMap {
		fmt.Printf("%s\t%0.2f\t%0.2f\t%0.2f\n", key, citiesResultMap[key].min, citiesResultMap[key].max, citiesResultMap[key].average)
	}

	fmt.Printf("\nQuantity of cities: %d\n", len(citiesResultMap))
}

func readLines(scanner *bufio.Scanner, linesChannel chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	linesBatch := make([]string, 0, batchSize)
	for scanner.Scan() {
		newLine := scanner.Text()
		linesBatch = append(linesBatch, newLine)
		if len(linesBatch) == batchSize {
			linesChannel <- linesBatch
			linesBatch = linesBatch[:0]
		}
	}
	if len(linesBatch) > 0 {
		linesChannel <- linesBatch
	}
	linesChannel <- linesBatch
	close(linesChannel)
}

func parseLines(citiesMap map[string]cityAggregate, linesChannel chan []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for linesBatch := range linesChannel {
		for _, line := range linesBatch {
			cityName, temperatureString, _ := strings.Cut(line, ";")
			temperature, _ := parseFloat(temperatureString)

			currentCity := citiesMap[cityName]
			currentCity.count++
			currentCity.sum += temperature
			if temperature > currentCity.max {
				currentCity.max = temperature
			}
			if temperature < currentCity.min {
				currentCity.min = temperature
			}

			citiesMap[cityName] = currentCity
		}
	}
}

func parseFloat(temperatureStr string) (float64, error) {
	temperatureIntStr, fractionStr, ok := strings.Cut(temperatureStr, ".")
	if !ok {
		return 0, fmt.Errorf("error splitting number string")
	}
	temperatureInt, err := strconv.Atoi(temperatureIntStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing, %v", err)
	}
	fraction, err := strconv.Atoi(fractionStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing, %v", err)
	}

	sign := +1
	if temperatureInt < 0 {
		sign = -1
		temperatureInt = -temperatureInt
	}
	hundredths := fraction
	decimal := fraction / 10
	hundredthsRemaining := hundredths % 10

	if sign == 1 && hundredthsRemaining > 0 {
		decimal++
	}

	temperature := float64((10*temperatureInt+decimal)*sign) / 10
	return temperature, nil
}

func approxToDecimal(input float64) float64 {
	output := math.Round(float64((input * 10)) / 10)
	return output
}
