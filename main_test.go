package main

import (
	"testing"
	"time"
)

func Test_minmax(t *testing.T) {
	type args struct {
		data []Data
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 float64
	}{
		{"Normal array", args{[]Data{
			{time.Now(), 1},
			{time.Now(), 2},
			{time.Now(), 3},
		}}, 1, 3},
		{"Inverted array", args{[]Data{
			{time.Now(), 3},
			{time.Now(), 2},
			{time.Now(), 1},
		}}, 1, 3},
		{"All same value", args{[]Data{
			{time.Now(), 1},
			{time.Now(), 1},
			{time.Now(), 1},
		}}, 1, 1},
		{"Only one value", args{[]Data{
			{time.Now(), 2},
		}}, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := minmax(tt.args.data)
			if got != tt.want {
				t.Errorf("minmax() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("minmax() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_calculateEpsilon(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"Normal range", args{0.0, 35.0}, 10},
		{"Small range", args{0.0, 0.01}, 0.001},
		{"Extremely small range", args{-0.99, 0.001}, 0.1},
		{"Non-existing range", args{0.0, 0.0}, 0.0001},
		{"Inverted range", args{-18, 1}, 1},
		{"Extremely small values", args{-0.00002, 0.00001}, 0.0001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateEpsilon(tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("calculateEpsilon() = %v, want %v", got, tt.want)
			}
		})
	}
}
