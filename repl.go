package main

import "strings"

func CleanInput(text string) []string {
	fields := strings.Fields(strings.ToLower(text))
	return fields
}
