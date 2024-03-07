package prometheus_handler

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
)

const ( // constants defined for Metrics
	HISTOGRAM = 0
	COUNTER   = 1
	GAUGE     = 2
	UNTYPE    = 3
)

type HandlerStructure []struct {
	MType    int               // Field to store the Metric Type
	MName    string            // Field to store Metric Name
	LabelMap map[string]string // Map to store Metric Labels
	MValue   interface{}       // Field to store Metric Value
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

func makePromCounter(label string, count string, labelData map[string]string) string {
	var entry, labelStr string
	output := fmt.Sprintf(`
# HELP %s counter output
# TYPE %s counter`, label, label)
	for k, v := range labelData {
		labelStr += fmt.Sprintf(`%s="%s", `, regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(k, ""), v)
	}
	if len(labelData) > 0 {
		entry = fmt.Sprintf("%s\n%s{%s} %s", output, label, labelStr, count)
	} else {
		entry = fmt.Sprintf("%s\n%s %s", output, label, count)
	}
	return entry + "\n"
}

func makePromHistogram(label string, histogram map[string]int, labelData map[string]string) string {
	label = regexp.MustCompile(`[^a-zA-Z0-9_ ]+`).ReplaceAllString(label, "")
	var output, labelStr string
	output = fmt.Sprintf(`
# HELP %s histogram output
# TYPE %s histogram`, label, label)
	//Sort bound
	var bounds []string
	for bound := range histogram {
		bounds = append(bounds, bound)
	}
	sort.Strings(bounds)
	for k, v := range labelData {
		labelStr += fmt.Sprintf(`%s="%s", `, regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(k, ""), v)
	}

	for _, bound := range bounds {
		var entry string
		if len(labelData) > 0 {
			entry = fmt.Sprintf(`%s_bucket{%s le="%s"} %d`, label, labelStr, bound, histogram[bound])
		} else {
			entry = fmt.Sprintf(`%s_bucket{le="%s"} %d`, label, bound, histogram[bound])
		}
		output += "\n" + entry
	}

	var entry1, entry2 string
	if len(labelData) > 0 {
		entry1 = fmt.Sprintf(`%s_sum{%s} %d`, label, labelStr, histogram["+inf"])
		entry2 = fmt.Sprintf(`%s_count{%s} %d`, label, labelStr, histogram["+inf"])
	} else {
		entry1 = fmt.Sprintf(`%s_sum %d`, label, histogram["+inf"])
		entry2 = fmt.Sprintf(`%s_count %d`, label, histogram["+inf"])
	}
	output += "\n" + entry1 + "\n" + entry2 + "\n"
	return output
}

func makePromGauge(label string, value string, labelData map[string]string) string {
	var entry, labelStr string
	output := fmt.Sprintf(`
# HELP %s gauge output
# TYPE %s gauge`, label, label)
	for k, v := range labelData {
		labelStr += fmt.Sprintf(`%s="%s", `, regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(k, ""), v)
	}
	if len(labelData) > 0 {
		entry = fmt.Sprintf("%s\n%s{%s} %s", output, label, labelStr, value)
	} else {
		entry = fmt.Sprintf("%s\n%s %s", output, label, value)
	}
	return entry + "\n"
}

func GenericPromDataParser(structure HandlerStructure) string {
	var data string
	for i := 0; i < len(structure); i++ {
		{
			var op string
			promType := structure[i].MType                     // Gets the MetricType from HandlerStructure Structure
			promLabel := structure[i].MName                    // Gets the MetricName from HandlerStructure Structure
			fieldValue := reflect.ValueOf(structure[i].MValue) // Gets the MetricValue from HandlerStructure Structur
			switch promType {
			case HISTOGRAM:
				histogram := parseHistogram(fieldValue)
				op += makePromHistogram(promLabel, histogram, structure[i].LabelMap)
			case COUNTER:
				count := parseCounter(fieldValue)
				op += makePromCounter(promLabel, count, structure[i].LabelMap)
			case GAUGE:
				value := parseGauge(fieldValue)
				op += makePromGauge(promLabel, value, structure[i].LabelMap)
			case UNTYPE:
				op += makePromUntype(promLabel, fieldValue)
			}
			data = data + op
		}
	}
	return data
}
