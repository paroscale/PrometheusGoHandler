package prometheus_handler

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	//"strconv"
	//"strings"
)

const (
	HISTOGRAM = 0
	COUNTER   = 1
	GAUGE     = 2
	UNTYPE    = 3
)

type HandlerStructure []struct {
	MType   int
	MName   string
	MLName  string
	MLValue string
	MValue  interface{}
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
	if len(histoMap) > 0 {
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
# HELP %s counter output
# TYPE %s counter
`, label, label)
	entry := fmt.Sprintf(`%s%s %s`, output, label, count)
	return entry + "\n"
}

func makePromHistogram(label string, histogram map[string]int, MLName string, MLValue string) string {
	output := fmt.Sprintf(`
# HELP %s histogram output
# TYPE %s histogram`, label, label)

	//Sort bound
	var bounds []string
	for bound, _ := range histogram {
		bounds = append(bounds, bound)
	}
	sort.Strings(bounds)
	MLName = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(MLName, "")
	if MLValue == "1" {
		MLValue = MLName
	}
	for _, bound := range bounds {
		entry := fmt.Sprintf(`%s_bucket{le="%s", %s="%s"} %d`, label, bound, strings.TrimRight(MLName, "bt"), MLValue, histogram[bound])
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
	for i := 0; i < len(structure); i++ {
		{
			var op string
			fieldValue := reflect.ValueOf(structure[i].MValue)
			promType := structure[i].MType
			promLabel := structure[i].MName
			switch promType {
			case HISTOGRAM:
				histogram := parseHistogram(fieldValue)
				op += makePromHistogram(promLabel, histogram, structure[i].MLName, structure[i].MLValue)
			case COUNTER:
				count := parseCounter(fieldValue)
				op += makePromCounter(promLabel, count)
			case GAUGE:
				value := parseGauge(fieldValue)
				op += makePromGauge(promLabel, value)
			case UNTYPE:
				op += makePromUntype(promLabel, fieldValue)
			}
			data = data + op
		}
	}
	return data
}
