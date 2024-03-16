package main

import (
	"bufio"
	"fmt"
	"math/rand"
  "runtime/pprof"
	"os"
)

const qtyBatches = 1 // 100,000 batches to generate 1 billion rows, assuming each batch processes 10,000 rows

type cityMeasure struct {
	name string
	data float32
}

func main(){
	f, err := os.Create("cpu.pprof")
  if err != nil {
    panic(err)
  }
  pprof.StartCPUProfile(f)
  defer pprof.StopCPUProfile()


	filename := "data/cities-10k.txt"
	DumpCityMeasuresFile(filename, "output.csv")
}

// DumpCityMeasuresFile generates city measurements and writes them to a file in batches.
func DumpCityMeasuresFile(inputFile, outputFile string) error {
	cities, err := GetCities(inputFile)
	if err != nil {
		return fmt.Errorf("failed to get cities: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Loop through the number of batches required to reach the goal
	for i := 0; i < qtyBatches; i++ {
		for j := 0; j < len(cities); j++ {
			// Simulate getting a city measurement (you may need to adjust the logic based on your city list size)
			cityIndex := rand.Intn(len(cities)) // Randomly pick a city to simulate variability in data
			measure := cityMeasure{name: cities[cityIndex], data: getRandomTemperature()}

			// Write the city measure to the buffer
			if _, err := writer.WriteString(fmt.Sprintf("%s;%0.2f\n", measure.name, measure.data)); err != nil {
				return fmt.Errorf("failed to write city measure: %w", err)
			}
		}

		// Flush the buffer after each batch to write to file
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush writer: %w", err)
		}
	}

	return nil
}

func GetCities(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cities []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newToken := scanner.Text()
		cities = append(cities, newToken)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cities, err
}

func getRandomTemperature() float32 {
	return (2 * rand.Float32() * 99.99) - 99.99
}
