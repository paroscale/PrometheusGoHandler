package prometheus_handler

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"regexp"
	//"strconv"
	//"strings"
)

type HandlerStructure []struct {
	Structure interface{}
	Datatype  string
	Labelname string
}

func parseCounter(count reflect.Value) string {
	return fmt.Sprintf("%v", count)
}

func parseHistogram(histogram reflect.Value) map[string]int {
	inter := histogram.Interface()
	histoMap := inter.(map[string]int)
	totalObservation := 0
	for _, value := range histoMap {
		totalObservation += value
	}
	if len(histoMap)>0 {
		histoMap["+inf"] = totalObservation
	}

	return histoMap
}

func parseGauge(gauge reflect.Value) string {
	return fmt.Sprintf("%v", gauge)
}

func makePromUntype(label string, value reflect.Value) string {
	strip := fmt.Sprintf("%v", value)
	output := fmt.Sprintf(`
# UNtype output
%s %s`, label, strip)
	return output
}

func makePromCounter(label string, count string) string {
	output := fmt.Sprintf(`
# HELP %s outputr cell such that its distance to the nearest land cell is maximized, and return the distance. If no land or water exists in the grid, return -1.

The distance used in this problem is the Manhattan distance: the distance between two cells (x0, y0) and (x1, y1) is |x0 - x1| + |y0 - y1|.

 
# TYPE %s counter\n`, label, label)
	entry := fmt.Sprintf(`%s%s %s`, output, label, count)
	return entry + "\n"
}

func makePromHistogram(label string, histogram map[string]int, datatype string, labelname string) string {
	output := fmt.Sprintf(`
# HELP %s histogram output
# TYPE %s histogram`, label, label)

	//Sort bound
	var bounds []string
	for bound, _ := range histogram {
		bounds = append(bounds, bound)
	}
	sort.Strings(bounds)
	labelname = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(labelname, "")
	if datatype=="1"{
		datatype=labelname
	}
	for _, bound := range bounds {
		entry := fmt.Sprintf(`%s_bucket{le="%s", %s="%s"} %d`, label, bound, strings.TrimRight(labelname, "bt"), datatype, histogram[bound])
		output += "\n" + entry
	}
	entry := fmt.Sprintf("%s_count %d", label, histogram["+inf"])
	output += "\n\n" + entry + "\n"
	return output
}

func makePromGauge(label string, value string) string {
	output := fmt.Sprintf(`
# HELP %s gauge output
# TYPE %s gauge 
`, label, label)
	entry := fmt.Sprintf(`%s%s %s`, output, label, value)
	return entry + "\n"
}

func GenericPromDataParser(structure HandlerStructure) string {
	var data string
	for j := 0; j < len(structure); j++ {
		{
			var op string
			//Reflect the struct
			typeExtract := reflect.TypeOf(structure[j].Structure)
			valueExtract := reflect.ValueOf(structure[j].Structure)

			//Iterating over fields and compress their values
			for i := 0; i < typeExtract.NumField(); i++ {
				fieldType := typeExtract.Field(i)
				fieldValue := valueExtract.Field(i)

				promType := fieldType.Tag.Get("type")
				promLabel := fieldType.Tag.Get("metric")
				switch promType {
				case "histogram":
					histogram := parseHistogram(fieldValue)
					op += makePromHistogram(promLabel, histogram, structure[j].Datatype, structure[j].Labelname)
				case "counter":
					count := parseCounter(fieldValue)
					op += makePromCounter(promLabel, count)
				case "gauge":
					value := parseGauge(fieldValue)
					op += makePromGauge(promLabel, value)
				case "untype":
					op += makePromUntype(promLabel, fieldValue)
				}
			}
			data = data + op
		}
	}
	return data
}
