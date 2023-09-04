package main

import (
	"fmt"
	prometheus_handler "handler/prometheus"
	"time"
)

var addToHandler prometheus_handler.HandlerStructure

func main() {

	for i := 0; i < 5; i++ {
		addToHandler = append(addToHandler, struct {
			MType   int
			MName   string
			MLName  string
			MLValue string
			MValue  interface{}
		}{MType: prometheus_handler.GAUGE, MName: "Field1", MValue: i, MLName: "label", MLValue: "labelValue"})
	}
	time.Sleep(2 * time.Second)
	for i := 5; i > 0; i-- {
		addToHandler = append(addToHandler, struct {
			MType   int
			MName   string
			MLName  string
			MLValue string
			MValue  interface{}
		}{MType: prometheus_handler.GAUGE, MName: "Field1", MValue: i})
	}
	result := prometheus_handler.GenericPromDataParser(addToHandler)
	fmt.Println(result)
	time.Sleep(2 * time.Second)
}
