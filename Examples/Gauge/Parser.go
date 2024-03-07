package main

import (
	"fmt"
	prometheus_handler "handler/prometheus"
	"time"
)

var addToHandler prometheus_handler.HandlerStructure

func main() {
	lm := make(map[string]string)
	lm["Label1"] = "value1"
	lm["Label2"] = "value2"
	for i := 0; i < 5; i++ {
		addToHandler = append(addToHandler, struct {
			MType    int
			MName    string
			LabelMap map[string]string
			MValue   interface{}
		}{MType: prometheus_handler.GAUGE, MName: "Field1", MValue: i, LabelMap: lm})
	}
	time.Sleep(2 * time.Second)
	for i := 5; i > 0; i-- {
		addToHandler = append(addToHandler, struct {
			MType    int
			MName    string
			LabelMap map[string]string
			MValue   interface{}
		}{MType: prometheus_handler.GAUGE, MName: "Field1", MValue: i, LabelMap: lm})
	}
	result := prometheus_handler.GenericPromDataParser(addToHandler)
	fmt.Println(result)
	time.Sleep(2 * time.Second)
}
