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

var Result string
var addToHandler prometheus_handler.HandlerStructure

func main() {
	var Hist Histogram
	lm := make(map[string]string)
	lm["Label1"] = "value1"
	lm["Label2"] = "value2"

	Hist.Num0 = 1
	Hist.Num2 = 2
	Hist.Num4 = 3
	var histmap = make(map[string]int)
	histmap[strconv.Itoa(Hist.Num0)] = Hist.Num2
	histmap[strconv.Itoa(Hist.Num2)] = Hist.Num4
	histmap[strconv.Itoa(Hist.Num4)] = Hist.Num0
	fmt.Println(histmap)
	addToHandler = append(addToHandler, struct {
		MType    int
		MName    string
		MLName   string
		MLValue  string
		LabelMap map[string]string
		MValue   interface{}
	}{MType: prometheus_handler.HISTOGRAM, MName: "hist_field1", MValue: histmap, LabelMap: lm})
	result := prometheus_handler.GenericPromDataParser(addToHandler)
	fmt.Println(result)
}
