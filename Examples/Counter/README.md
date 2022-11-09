
# Import Counter using the Handler

To Import a Counter using the Promethues Handler, Follow the Steps Below.

- As the Export Type will be Counter, we need to Create a Export Structure to have type as `int`.

```go
type Export struct {
 Field1 map[string]int `type:"histogram" metric:"Field1"`
 Field2 int `type:"counter" metric:"Field2"`
 Field3 int `type:"gauge" metric:"Field3"`
}
```

- Pass the `var Counter int` to the Export Structure and then pass the Structure to `func GenericPromDataParser(structure interface{}) string`

- In `func GenericPromDataParser`

	- First the Function will look into the `type` field to determine what type of Parser Function needs to be called. In our case as it'll be Counter, `func parseCounter` will be called.
	- It'll convert the Counter to the required type and return back to the `GenericPromDataParser`.
	- Then `func makePromCounter` will be called with the Counter and `metric` field value.
	- This Function will create the Strings which Promethues can Parse easily and returns the Output which could be posted on the Port where Promethues reads the data from.

## Example Output

- For the Example Program we ran a Counter which increments by 1 every 2 Seconds and call `func GenericPromDataParser`

- This is the Output which will be given back from `func GenericPromDataParser`

```
# HELP Field1 output
# TYPE Field1 counter
Field1 0


# HELP Field1 output
# TYPE Field1 counter
Field1 1


# HELP Field1 output
# TYPE Field1 counter
Field1 2

.
.
.
.

# HELP Field1 output
# TYPE Field1 counter
Field1 10
```
