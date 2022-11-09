package main

import (
	"fmt"
	prometheus_handler "handler/prometheus"
	"time"
)

type Export struct {
	Field1 int `type:"counter" metric:"Field1"`
}

func main() {
	var ht Export
	for i := 0; i < 10; i++ {
		ht.Field1 = i
		result := prometheus_handler.GenericPromDataParser(ht)
		fmt.Println(result)
		time.Sleep(2 * time.Second)
	}
}
