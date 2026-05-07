package main

import (
	"math"
	"testing"
	"time"
)

func TestMinMax(t *testing.T) {
	tests := []struct {
		name        string
		data        []Data
		wantMin     float64
		wantMax     float64
		wantErr     bool
	}{
		{
			name: "Normal array",
			data: []Data{
				{time.Now(), 1},
				{time.Now(), 2},
				{time.Now(), 3},
			},
			wantMin: 1,
			wantMax: 3,
		},
		{
			name: "Inverted array",
			data: []Data{
				{time.Now(), 3},
				{time.Now(), 2},
				{time.Now(), 1},
			},
			wantMin: 1,
			wantMax: 3,
		},
		{
			name: "All same value",
			data: []Data{
				{time.Now(), 1},
				{time.Now(), 1},
				{time.Now(), 1},
			},
			wantMin: 1,
			wantMax: 1,
		},
		{
			name: "Only one value",
			data: []Data{
				{time.Now(), 2},
			},
			wantMin: 2,
			wantMax: 2,
		},
		{
			name:    "Empty array",
			data:    []Data{},
			wantErr: true,
		},
		{
			name: "Array with NaN",
			data: []Data{
				{time.Now(), 1},
				{time.Now(), math.NaN()},
				{time.Now(), 3},
			},
			wantMin: 1,
			wantMax: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax, err := MinMax(tt.data)
			if tt.wantErr && err == nil {
				t.Errorf("MinMax() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("MinMax() unexpected error: %v", err)
				return
			}
			if !tt.wantErr {
				if gotMin != tt.wantMin {
					t.Errorf("MinMax() gotMin = %v, want %v", gotMin, tt.wantMin)
				}
				if gotMax != tt.wantMax {
					t.Errorf("MinMax() gotMax = %v, want %v", gotMax, tt.wantMax)
				}
			}
		})
	}
}

func TestCalculateEpsilon(t *testing.T) {
	tests := []struct {
		name string
		min  float64
		max  float64
		want float64
	}{
		{"Normal range", 0.0, 30.0, 2.0},
		{"Small range", 0.0, 0.01, (0.01 / 15)},
		{"Extremely small range", -0.99, 0.001, (0.991 / 15)},
		{"Non-existing range", 0.0, 0.0, 0.0001},
		{"Inverted range", 16, 1, 1.0},
		{"Extremely small values", -0.00002, 0.00001, 0.0001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateEpsilon(tt.min, tt.max)
			// Use approximate comparison for floating point
			if math.Abs(got-tt.want) > 0.000001 {
				t.Errorf("CalculateEpsilon() = %v, want %v (diff: %v)", got, tt.want, math.Abs(got-tt.want))
			}
		})
	}
}

func TestConvertDataArray(t *testing.T) {
	jsonData := []JsonData{
		{Rank: 1, DateTime: "2023-01-01T12:00:00+00:00", PriceNoTax: 10.0, PriceWithTax: 12.0},
		{Rank: 2, DateTime: "2023-01-01T13:00:00+00:00", PriceNoTax: 20.0, PriceWithTax: 24.0},
	}

	got, err := ConvertDataArray(jsonData)
	if err != nil {
		t.Fatalf("ConvertDataArray() unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("ConvertDataArray() expected 2 elements, got %d", len(got))
	}

	if got[0].Price != 12.0 {
		t.Errorf("ConvertDataArray() first price = %v, want 12.0", got[0].Price)
	}
	if got[1].Price != 24.0 {
		t.Errorf("ConvertDataArray() second price = %v, want 24.0", got[1].Price)
	}
}
