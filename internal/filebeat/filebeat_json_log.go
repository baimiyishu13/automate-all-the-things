package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a string argument")
		return
	}

	input := os.Args[1]
	lines := strings.Split(input, "\n")

	logInfo := map[string][]string{
		"json_log": lines,
	}

	jsonData, err := json.MarshalIndent(logInfo, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	if err := os.WriteFile("../../test/json_log.json", jsonData, 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("Successfully wrote to json_log.json")
}
