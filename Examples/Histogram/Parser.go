package main

import (
	"fmt"
	prometheus_handler "handler/prometheus"
	"strconv"
)

type Histogram struct {
	Num0 int `le:0`
	Num2 int `le:2`
	Num4 int `le:4`
}

type Export struct {
	Field1 map[string]int `type:"histogram" metric:"Field1"`
}

var Result string

func main() {
	var Hist Histogram
	Hist.Num0 = 1
	Hist.Num2 = 2
	Hist.Num4 = 3
	var histmap = make(map[string]int)
	histmap[strconv.Itoa(Hist.Num0)] = Hist.Num2
	histmap[strconv.Itoa(Hist.Num2)] = Hist.Num4
	histmap[strconv.Itoa(Hist.Num4)] = Hist.Num0
	fmt.Println(histmap)
	var ht Export
	ht.Field1 = histmap
	result := prometheus_handler.GenericPromDataParser(ht)
	fmt.Println(result)
}
