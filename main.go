package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

const URL = "https://api.spot-hinta.fi/TodayAndDayForward"
const TS_LAYOUT = "2006-01-02T15:04:05-07:00"
const HTTP_TIMEOUT = 10 * time.Second
const GRAPH_STEP = 15.0
const MAX_DISPLAY_POINTS = 80

// Data is now 15-minute resolution (4 points per hour)
// Downsample to hourly by keeping every 4th point
const DOWNAMPLE_STEP = 4

// Hour markers for graph display
const (
	HourMarker6  = 6
	HourMarker12 = 12
)

type JsonData struct {
	Rank         int     `json:"Rank"`
	DateTime     string  `json:"DateTime"`
	PriceNoTax   float64 `json:"PriceNoTax"`
	PriceWithTax float64 `json:"PriceWithTax"`
}

type Data struct {
	Time  time.Time
	Price float64
}

// ErrEmptyData is returned when no data is available
var ErrEmptyData = errors.New("no data available")

func GetDataArray(url string) ([]JsonData, error) {
	client := http.Client{Timeout: HTTP_TIMEOUT}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	var dataArray []JsonData
	err = json.Unmarshal(data, &dataArray)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return dataArray, nil
}

func ConvertDataArray(jsonData []JsonData) ([]Data, error) {
	var result []Data

	for _, value := range jsonData {
		dateTime, err := time.Parse(TS_LAYOUT, value.DateTime)
		if err != nil {
			return nil, fmt.Errorf("failed to convert timestamp: %w", err)
		}
		price := value.PriceWithTax
		result = append(result, Data{Time: dateTime, Price: price})
	}

	return result, nil
}

// DownsampleData reduces the data points to fit within display constraints
// For 15-minute data (4 points per hour), we keep every 4th point to get hourly resolution
func DownsampleData(data []Data, step int) []Data {
	if len(data) == 0 {
		return nil
	}

	var result []Data
	for i := 0; i < len(data); i += step {
		result = append(result, data[i])
	}

	return result
}

func MinMax(data []Data) (float64, float64, error) {
	if len(data) == 0 {
		return 0, 0, ErrEmptyData
	}

	min := data[0].Price
	max := min

	for _, d := range data {
		if !math.IsNaN(d.Price) {
			if max < d.Price {
				max = d.Price
			}
			if min > d.Price {
				min = d.Price
			}
		}
	}

	return min, max, nil
}

func CalculateEpsilon(minimum float64, maximum float64) float64 {
	minValue := math.Min(minimum, maximum)
	maxValue := math.Max(minimum, maximum)
	minValue = math.Round(minValue*1000) / 1000
	maxValue = math.Round(maxValue*1000) / 1000
	delta := maxValue - minValue

	if delta == 0.0 {
		return 0.0001
	}

	return delta / GRAPH_STEP
}

func GenerateGraph(dataArray []Data, min float64, max float64, epsilon float64) {
	for i := max; i >= min; i -= epsilon {
		fmt.Printf("\n%7.4f | ", i)
		for _, data := range dataArray {
			if math.IsNaN(data.Price) {
				fmt.Print(" ")
			} else if i <= data.Price {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
	}
	fmt.Print("\n  €/kWh ")
	for _, data := range dataArray {
		if data.Time.Hour()%HourMarker6 == 0 {
			fmt.Print("''|'''")
		}
	}
	fmt.Print("\n      ")
	for _, data := range dataArray {
		if data.Time.Hour()%HourMarker12 == 0 {
			fmt.Printf("    ^       ")
		}
	}
	fmt.Print("\n      ")
	for _, data := range dataArray {
		if data.Time.Hour()%HourMarker12 == 0 {
			fmt.Printf("  %02d.%02d     ", data.Time.Day(), data.Time.Month())
		}
	}
	fmt.Print("\n      ")
	for _, data := range dataArray {
		if data.Time.Hour()%HourMarker12 == 0 {
			fmt.Printf("  %02d:%02d     ", data.Time.Hour(), data.Time.Minute())
		}
	}
	fmt.Printf("\n\nMin: %.4f, Max: %.4f\n", min, max)
}

func main() {
	jsonData, err := GetDataArray(URL)
	if err != nil {
		log.Fatalf("Failed to get data array: %v", err)
	}

	dataArray, err := ConvertDataArray(jsonData)
	if err != nil {
		log.Fatalf("Failed to convert data array: %v", err)
	}

	// Calculate min/max from ALL data points (before downsampling)
	// This ensures the Y-axis range reflects the actual price range
	min, max, err := MinMax(dataArray)
	if err != nil {
		log.Fatalf("Failed to calculate min/max: %v", err)
	}

	// Downsample 15-minute data to hourly (keep every 4th point)
	// This ensures the X-axis fits within 80 columns
	dataArray = DownsampleData(dataArray, DOWNAMPLE_STEP)

	// Add marker for end of graph
	last := dataArray[len(dataArray)-1].Time
	last = last.Add(time.Hour)
	dataArray = append(dataArray, Data{Time: last, Price: math.NaN()})

	epsilon := CalculateEpsilon(min, max)
	GenerateGraph(dataArray, min, max, epsilon)
}
