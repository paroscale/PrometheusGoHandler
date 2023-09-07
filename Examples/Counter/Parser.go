package main

import (
	"fmt"
	prometheus_handler "handler/prometheus"
	"time"
)

func main() {
	var addToHandler prometheus_handler.HandlerStructure
	for i := 0; i < 10; i++ {
		addToHandler = append(addToHandler, struct {
			MType   int
			MName   string
			MLName  string
			MLValue string
			MValue  interface{}
		}{MType: prometheus_handler.COUNTER, MName: "Field1", MValue: i, MLName: "counter", MLValue: "test counter"})
		result := prometheus_handler.GenericPromDataParser(addToHandler)
		addToHandler = nil
		fmt.Println(result)
		time.Sleep(2 * time.Second)
	}
}
