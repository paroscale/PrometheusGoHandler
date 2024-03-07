
# Import Counter using the Handler

To Import a Counter using the Promethues Handler, Follow the Steps Below.

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
		LabelMap map[string]string
		MValue  interface{}
	}{MType: prometheus_handler.COUNTER, MName: "Field1", MValue: 20, LabelMap: labelmap})
```

- Pass the HandlerStructure and Call the Function `func GenericPromDataParser(structure HandleStructure) string` and set the `MTYPE` field to `COUNTER`
	
- In `func GenericPromDataParser`

	- First the Function will look into the `MTYPE` field to determine what type of Parser Function needs to be called. In our case as it'll be Counter, `func parseCounter` will be called.
	- It'll convert the Counter to the required type and return back to the `GenericPromDataParser`.
	- Then `func makePromCounter` will be called with the Counter and `metric` field value.
	- This Function will create the Strings which Promethues can Parse easily and returns the Output which could be posted on the Port where Promethues reads the data from.

## Example Output

- For the Example Program we ran a Counter which increments by 1 every 2 Seconds and call `func GenericPromDataParser`

- This is the Output which will be given back from `func GenericPromDataParser`

## Without Label
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

## With Label

```
# HELP Field1 counter output
# TYPE Field1 counter
Field1{Label1="value1", Label2="value2", } 0

# HELP Field1 counter output
# TYPE Field1 counter
Field1{Label1="value1", Label2="value2", } 1

# HELP Field1 counter output
# TYPE Field1 counter
Field1{Label1="value1", Label2="value2", } 2
```