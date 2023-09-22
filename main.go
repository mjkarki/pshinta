package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

const URL = "https://api.spot-hinta.fi/TodayAndDayForward"
const TS_LAYOUT = "2006-01-02T15:04:05-07:00"

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

func getDataArray(url string) []JsonData {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to get data:", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read data:", err)
	}

	var dataArray []JsonData
	err = json.Unmarshal(data, &dataArray)
	if err != nil {
		log.Fatal("Failed to parse JSON:", err)
	}

	return dataArray
}

func convertDataArray(jsonData []JsonData) []Data {
	var result []Data

	for _, value := range jsonData {
		dt, err := time.Parse(TS_LAYOUT, value.DateTime)
		if err != nil {
			log.Fatal("Failed to convert timestamp:", err)
		}
		pr := value.PriceWithTax
		result = append(result, Data{Time: dt, Price: pr})
	}

	return result
}

func minmax(data []Data) (float64, float64) {
	var min float64 = data[0].Price
	var max float64 = min

	for _, d := range data {
		if max < d.Price {
			max = d.Price
		}
		if min > d.Price {
			min = d.Price
		}
	}

	return min, max
}

func calculateEpsilon(min float64, max float64) float64 {
	max = math.Round(max*1000) / 1000
	min = math.Round(min*1000) / 1000
	if max-min == 0.0 {
		return 0.0001
	}

	var exp int = int(math.Round(math.Log10(math.Abs(max-min)))) - 1
	var epsilon float64 = math.Pow10(exp)

	return epsilon
}

func generateGraph(dataArray []Data, min float64, max float64, epsilon float64) {
	for i := (max + epsilon); i >= min; i -= epsilon {
		fmt.Printf("\n%7.4f | ", i)
		for _, data := range dataArray {
			if i < (data.Price + epsilon) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
	}
	fmt.Print("\n€/kWh ")
	for _, data := range dataArray {
		if data.Time.Hour()%6 == 0 {
			fmt.Print("''|'''")
		}
	}
	fmt.Print("\n    ")
	for _, data := range dataArray {
		if data.Time.Hour()%12 == 0 {
			fmt.Printf("    ^       ")
		}
	}
	fmt.Print("\n    ")
	for _, data := range dataArray {
		if data.Time.Hour()%12 == 0 {
			fmt.Printf("  %02d.%02d     ", data.Time.Day(), data.Time.Month())
		}
	}
	fmt.Print("\n    ")
	for _, data := range dataArray {
		if data.Time.Hour()%12 == 0 {
			fmt.Printf("  %02d:%02d     ", data.Time.Hour(), data.Time.Minute())
		}
	}
	fmt.Println()
}

func main() {
	dataArray := convertDataArray(getDataArray(URL))
	last := dataArray[len(dataArray)-1].Time
	last = last.Add(time.Hour)
	dataArray = append(dataArray, Data{Time: last, Price: math.NaN()})

	min, max := minmax(dataArray)
	epsilon := calculateEpsilon(min, max)

	generateGraph(dataArray, min, max, epsilon)
}
