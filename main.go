package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cityAggregate struct {
	min   float32
	max   float32
	sum   float32
	count int
}

type cityResult struct {
	min     float32
	max     float32
	average float32
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
		currentCity.sum += float32(temperature)
		if temperature > float64(currentCity.max) {
			currentCity.max = float32(temperature)
		}
		if temperature < float64(currentCity.min) {
			currentCity.min = float32(temperature)
		}

		citiesMap[cityName] = currentCity
	}

	citiesResultMap := make(map[string]cityResult)
	for key := range citiesMap {
		averageTemperature := citiesMap[key].sum / float32(citiesMap[key].count)
		currentCityResult := cityResult{min: citiesMap[key].min, max: citiesMap[key].max, average: averageTemperature}

		citiesResultMap[key] = currentCityResult
	}

	for key := range citiesResultMap {
		fmt.Printf("key: %s\tvalue: %+v\n", key, citiesMap[key])
	}

	fmt.Printf("\nQuantity of cities: %d\n", len(citiesResultMap))
}
