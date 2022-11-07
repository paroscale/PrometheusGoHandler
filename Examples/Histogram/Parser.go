package main

import (
	"encoding/json"
	"fmt"
	prometheus_handler "handler/prometheus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Bpftrace struct {
	Data struct {
		Usecs []struct {
			Count int `json:"count"`
			Max   int `json:"max"`
			Min   int `json:"min"`
		} `json:"@usecs"`
	} `json:"data"`
	Type string `json:"type"`
}

type BpftraceHistogram struct {
	Hist map[string]int `type:"histogram" metric:"hist"`
}

var OutputFilename = "Histogram_Example.json"
var Result string

func main() {
	fp, err := os.Open(OutputFilename)
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()
	bpftraceData, err := ioutil.ReadAll(fp)
	if err != nil {
		fmt.Println(err)
	}
	if len(bpftraceData) >= 1 {
		line := strings.Split(string(bpftraceData), "\n")
		var dt Bpftrace
		err = json.Unmarshal([]byte(line[0]), &dt)
		if err != nil {
			fmt.Println(err)
		}
		var histmap = make(map[string]int)
		for i := 0; i < len(dt.Data.Usecs); i++ {
			histmap[strconv.Itoa(dt.Data.Usecs[i].Max+1)] = dt.Data.Usecs[i].Count
		}
		fmt.Println(histmap)
		var ht BpftraceHistogram
		ht.Hist = histmap
		Result = prometheus_handler.GenericPromDataParser(ht)
	}
}
