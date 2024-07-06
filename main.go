package main

import (
	"fmt"
	"log"

	"github.com/mukezhz/json-go-struct/gojson"
)

func main() {
	jsonPath := "data.json"
	filePath := "fun_struct.go"
	packageName := "dtoss"

	jsonData, err := gojson.ReadJSONFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	gojson.Gen(jsonData, filePath, packageName)

	fmt.Println("DTO file written successfully")
}
