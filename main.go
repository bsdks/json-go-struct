package main

import (
	"fmt"
	"log"

	"github.com/bsdks/json-go-struct/gojson"
)

func main() {
	jsonPath := "data.json"
	filePath := "dto_struct.go"
	packageName := "dto"

	jsonData, err := gojson.ReadJSONFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	gojson.Gen(jsonData, filePath, packageName)

	fmt.Println("DTO file written successfully")
}
