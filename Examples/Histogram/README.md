
# Import Histogram using the Handler

To Import a Histogram using the Promethues Handler, Follow the Steps Below.

- Create a variable of type HandleStructure
```
	type HandlerStructure []struct {
		MType   int
		MName   string
		LabelMap map[string]string
		MValue  interface{}
	}
	var addToHandler prometheus_handler.HandlerStructure
```
- Append the data to HandlerStructure
```
	addToHandler = append(addToHandler, struct {
			MType   int
			MName   string
			MLName  string
			MLValue string
			MValue  interface{}
		}{MType: prometheus_handler.HISTOGRAM, MName: "Field1", MValue: histmap, LabelMap: labelmap})
```
- Pass the HandlerStructure and Call the Function `func GenericPromDataParser(structure HandleStructure) string` and set the `MTYPE` field to `HISTOGRAM`

- In `func GenericPromDataParser`

	- First the Function will look into the `MTYPE` field to determine what type of Parser Function needs to be called. In our case as it'll be Histogram, `func parseHistogram` will be called.
	- It'll go through the Map and Create a Histogram which will be returned back to the `GenericPromDataParser`.
	- Then `func makePromHistogram` will be called with the Histogram and `metric` field value.
	- This Function will create the Strings which Promethues can Parse easily and returns the Output which could be posted on the Port where Promethues reads the data from.

## Example Output

- For the Example Program we have used the below Map and Passed it into `func GenericPromDataParser`

```
map[1:2 2:3 3:1]
```
- This is the Output which will be given back from `func GenericPromDataParser`

```
# HELP hist_field1 histogram output
# TYPE hist_field1 histogram
hist_field1_bucket{Label1="value1", Label2="value2",  le="+inf"} 6
hist_field1_bucket{Label1="value1", Label2="value2",  le="1"} 2
hist_field1_bucket{Label1="value1", Label2="value2",  le="2"} 3
hist_field1_bucket{Label1="value1", Label2="value2",  le="3"} 1
hist_field1_sum{Label1="value1", Label2="value2", } 6
hist_field1_count{Label1="value1", Label2="value2", } 6

```
