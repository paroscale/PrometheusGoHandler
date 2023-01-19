package prometheus_handler

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	//"strconv"
	//"strings"
)

type HandlerStructure []struct {
	Structure interface{}
	DataType  string
	LabelName string
}

/*
type prometheusClient struct {
	HTTPoutput string
	Label map[string]string
}
*/

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
	histoMap["+inf"] = totalObservation

	return histoMap
	/*
		histoMap := make(map[string]string)
		totalObservation := 0
		for i := 0; i < histogram.NumField(); i++ {
			bucketBound := histogramType.Field(i)
			bucketValue := histogram.Field(i)
			count := int(bucketValue.Int())
			upperBound := strings.Split(bucketBound.Tag.Get("json"), ",")[0]
			histoMap[upperBound] = fmt.Sprintf("%v", bucketValue)

			totalObservation += count
		}
		histoMap["+inf"] = strconv.Itoa(int(totalObservation))
		return histoMap
	*/
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
# HELP %s output
# TYPE %s counter
`, label, label)
	entry := fmt.Sprintf(`%s%s %s`, output, label, count)
	return entry + "\n"
}

func makePromHistogram(label string, histogram map[string]int, dataType string,
	labelName string) string {
	output := fmt.Sprintf(`
# HELP %s histogram output
# TYPE %s histogram`, label, label)

	//Sort bound
	var bounds []string
	for bound, _ := range histogram {
		bounds = append(bounds, bound)
	}
	sort.Strings(bounds)

	for _, bound := range bounds {
		entry := fmt.Sprintf(`%s_bucket{le="%s", %s="%s"} %d`, label, bound,
			strings.TrimRight(labelName, ".bt"), dataType, histogram[bound])
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
			//Reflect the struct
			typeExtract := reflect.TypeOf(structure[i].Structure)
			valueExtract := reflect.ValueOf(structure[i].Structure)

			//Iterating over fields and compress their values
			for i := 0; i < typeExtract.NumField(); i++ {
				fieldType := typeExtract.Field(i)
				fieldValue := valueExtract.Field(i)

				promType := fieldType.Tag.Get("type")
				promLabel := fieldType.Tag.Get("metric")
				switch promType {
				case "histogram":
					histogram := parseHistogram(fieldValue)
					op += makePromHistogram(promLabel, histogram,
						structure[i].DataType, structure[i].LabelName)
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
