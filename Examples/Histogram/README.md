
# Import Histogram using the Handler

To Import a Histogram using the Promethues Handler, Follow the Steps Below.

- Create a Structure according to the Input Format that is going to be given.

```go
type <structName> struct {
 Field1 <dataType> `type:<metricType> metric:<metricName>`
}
```
For our Example, we've used:

```go
type Histogram struct {
	Num0 int `le:0`
	Num2 int `le:2`
	Num4 int `le:4`
}
```

- As the Export Type will be Histogram, we need to Create a Export Structure to have type as `map[string]int`.

```go
type Export struct {
 Field1 map[string]int `type:"histogram" metric:"Field1"`
 Field2 int `type:"counter" metric:"Field2"`
 Field3 int `type:"gauge" metric:"Field3"`
}
```

- Pass the Map to the Export Structure and Call the Function `func GenericPromDataParser(structure interface{}) string`

- In `func GenericPromDataParser`

	- First the Function will look into the `type` field to determine what type of Parser Function needs to be called. In our case as it'll be Histogram, `func parseHistogram` will be called.
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
# HELP Field1 histogram output
# TYPE Field1 histogram
Field1_bucket{le="+inf"} 6
Field1_bucket{le="1"} 2
Field1_bucket{le="2"} 3
Field1_bucket{le="3"} 1

Field1_count 6
```
