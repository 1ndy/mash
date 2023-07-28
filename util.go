package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func countIndent(node string) int {
	spaces := 0
	for _, c := range node {
		if c == ' ' {
			spaces += 1
		}
	}
	return spaces
}

func stripSpacesAndQuotes(node string) string {
	node_name_regex, _ := regexp.Compile(`\S+[^:]`)
	key := string(node_name_regex.Find([]byte(node)))
	quoteIndex := strings.Index(key, "\"")
	if quoteIndex != -1 {
		fmt.Fprintf(os.Stderr, "Found a quote in key %s. Mash does not currently support JSON-looking yaml. Please reformat and try again\n", key)
		os.Exit(1)
	}
	return key
}

func padString(str string, numSpaces int) string {
	spaces := ""
	for i := 0; i < numSpaces; i++ {
		spaces += " "
	}
	return spaces + str
}

func findMinimumIndent(keys []DocKey) int {
	min := keys[0].Spaces
	for _, key := range keys {
		if key.Spaces < min {
			min = key.Spaces
		}
	}
	return min
}

func findSpacingInterval(keys []DocKey) int {
	// extract number of spaces before each key
	spaces := make([]int, len(keys))
	for i, key := range keys {
		spaces[i] = key.Spaces
	}

	// take the derivative
	if len(spaces)-1 == 0 {
		return spaces[0]
	} else {
		derivative := make([]int, len(keys)-1)
		for i := 0; i < len(derivative); i++ {
			di := spaces[i+1] - spaces[i]
			derivative[i] = di
		}

		// see if any difference is not divided evenly
		firstVal := derivative[0]
		for i := 0; i < len(derivative); i++ {
			if derivative[i]%firstVal != 0 {
				fmt.Fprintln(os.Stderr, "Inconsistent spacing found in input, cannot determine hierarchy")
				os.Exit(1)
			}
		}
		return firstVal
	}
}
