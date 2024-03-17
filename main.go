package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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

func main() {
	/*
		iterate through file
			read line
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

	citiesMap := make(map[string]cityAggregate)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newLine := scanner.Text()

		cityName, temperatureStr, _ := strings.Cut(newLine, ";")
		temperature, _ := strconv.ParseFloat(temperatureStr, 32)

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

	citiesResultMap := make(map[string]cityResult)
	for key := range citiesMap {
		averageTemperature := math.Round(citiesMap[key].sum*10/float64(citiesMap[key].count)) / 10
		currentCityResult := cityResult{min: citiesMap[key].min, max: citiesMap[key].max, average: averageTemperature}

		citiesResultMap[key] = currentCityResult
	}

	for key := range citiesResultMap {
		fmt.Printf("key: %s\tvalue: %+v\n", key, citiesResultMap[key])
	}

	fmt.Printf("\nQuantity of cities: %d\n", len(citiesResultMap))
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
