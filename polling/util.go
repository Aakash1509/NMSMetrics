package main

import (
	"strconv"
	"strings"
)

func ParseResult(result string) interface{} {
	result = strings.TrimSpace(result)

	// If result is float
	if f, err := strconv.ParseFloat(result, 64); err == nil {
		return f
	}

	// If result is int
	if i, err := strconv.Atoi(result); err == nil {
		return i
	}

	// Else return as a string
	return result
}
