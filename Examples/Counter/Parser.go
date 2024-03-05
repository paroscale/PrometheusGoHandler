package main

import (
	"fmt"
	prometheus_handler "handler/prometheus"
	"time"
)

func main() {
	var addToHandler prometheus_handler.HandlerStructure
	lm := make(map[string]string)
	lm["Label1"] = "value1"
	lm["Label2"] = "value2"
	for i := 0; i < 10; i++ {
		addToHandler = append(addToHandler, struct {
			MType    int
			MName    string
			LabelMap map[string]string
			MValue   interface{}
		}{MType: prometheus_handler.COUNTER, MName: "Field1", MValue: i, LabelMap: lm})
		result := prometheus_handler.GenericPromDataParser(addToHandler)
		addToHandler = nil
		fmt.Println(result)
		time.Sleep(2 * time.Second)
	}
}
