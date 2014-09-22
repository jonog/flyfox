package main

import "strings"

func stringArrayToInterfaceArray(input []string) []interface{} {
	output := make([]interface{}, len(input))
	for i, v := range input {
		output[i] = interface{}(v)
	}
	return output
}

func resultsToJSONArray(strs []string) string {
	return "[" + strings.Join(strs, ",") + "]"
}

func typesAndResultsByTypeToJSON(types []string, results []string) string {
	var objects []string
	for idx, t := range types {
		objects = append(objects, "\""+t+"\":"+results[idx])
	}
	return "{" + strings.Join(objects, ",") + "}"
}
