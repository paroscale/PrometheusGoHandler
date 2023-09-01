package main

import (
	"fmt"
	prometheus_handler "handler/prometheus"
	"time"
)

type Export struct {
	Field1 int `type:"gauge" metric:"Field1"`
}

var addToHandler prometheus_handler.HandlerStructure

func main() {

	var ht Export
	for i := 0; i < 5; i++ {
		ht.Field1 = i
		addToHandler = append(addToHandler,struct{
			Structure interface{}
			Datatype  string
			Labelname string
		}{Structure: ht, Labelname: "Gauge"})
	}
	result := prometheus_handler.GenericPromDataParser(addToHandler)
	fmt.Println(result)
	time.Sleep(2 * time.Second)
	for i := 5; i > 0; i-- {
		ht.Field1 = i
		addToHandler = append(addToHandler, struct {
			Structure interface{}
			Datatype  string
			Labelname string
		}{Structure: ht, Labelname: "Gauge"})
	}
	result = prometheus_handler.GenericPromDataParser(addToHandler)
	fmt.Println(result)
	time.Sleep(2 * time.Second)
}
