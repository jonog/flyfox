package main

func stringArrayToInterfaceArray(input []string) []interface{} {
	output := make([]interface{}, len(input))
	for i, v := range input {
		output[i] = interface{}(v)
	}
	return output
}
