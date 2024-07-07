package gojson

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func generateStructFromJSON(jsonObj map[string]interface{}, structName string) (string, map[string]string) {
	var structFields []string
	nestedStructs := make(map[string]string)

	titleCaser := cases.Title(language.Und)

	for key, value := range jsonObj {
		fieldName := titleCaser.String(key)
		fieldType, nestedStructsFromValue := getFieldType(value, fieldName)
		structFields = append(structFields, fmt.Sprintf("%s %s `json:\"%s\"`", fieldName, fieldType, key))
		for nestedKey, nestedValue := range nestedStructsFromValue {
			nestedStructs[nestedKey] = nestedValue
		}
	}

	structDefinition := fmt.Sprintf("type %s struct {\n%s\n}", structName, strings.Join(structFields, "\n"))
	nestedStructs[structName] = structDefinition
	return structDefinition, nestedStructs
}

// Function to determine the Go type from a JSON value
func getFieldType(value interface{}, parentStructName string) (string, map[string]string) {
	nestedStructs := make(map[string]string)

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return "string", nestedStructs
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int", nestedStructs
	case reflect.Float32, reflect.Float64:
		return "float64", nestedStructs
	case reflect.Bool:
		return "bool", nestedStructs
	case reflect.Map:
		subStructName := parentStructName + "Nested"
		subStructDefinition, subNestedStructs := generateStructFromJSON(value.(map[string]interface{}), subStructName)
		nestedStructs[subStructName] = subStructDefinition
		for nestedKey, nestedValue := range subNestedStructs {
			nestedStructs[nestedKey] = nestedValue
		}
		return subStructName, nestedStructs
	case reflect.Slice:
		elemType, shouldReturn, returnValue, returnValue1 := handleSliceType(value, parentStructName, nestedStructs)
		if shouldReturn {
			return returnValue, returnValue1
		}
		return "[]" + elemType, nestedStructs
	default:
		return "interface{}", nestedStructs
	}
}

func handleSliceType(value interface{}, parentStructName string, nestedStructs map[string]string) (string, bool, string, map[string]string) {
	sliceElem := reflect.TypeOf(value).Elem()
	if len(value.([]interface{})) > 0 {
		firstElem := value.([]interface{})[0]
		if sliceElem.Kind() == reflect.Map {
			subStructName := parentStructName + "Item"
			subStructDefinition, subNestedStructs := generateStructFromJSON(firstElem.(map[string]interface{}), subStructName)
			nestedStructs[subStructName] = subStructDefinition
			for nestedKey, nestedValue := range subNestedStructs {
				nestedStructs[nestedKey] = nestedValue
			}
			return "", true, "[]" + subStructName, nestedStructs
		} else {
			elemType, subNestedStructs := getFieldType(firstElem, parentStructName)
			for nestedKey, nestedValue := range subNestedStructs {
				nestedStructs[nestedKey] = nestedValue
			}
			return "", true, "[]" + elemType, nestedStructs
		}
	}
	elemType, subNestedStructs := getFieldType(reflect.New(sliceElem).Elem().Interface(), parentStructName)
	for nestedKey, nestedValue := range subNestedStructs {
		nestedStructs[nestedKey] = nestedValue
	}
	return elemType, false, "", nil
}

func Gen(jsonData []byte, filePath string, packageName string) {
	var structDefinitions map[string]string
	if jsonData[0] == '[' {
		var jsonArray []interface{}
		err := json.Unmarshal(jsonData, &jsonArray)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON array: %v", err)
		}
		if len(jsonArray) > 0 {
			firstElem, ok := jsonArray[0].(map[string]interface{})
			if !ok {
				log.Fatalf("Error: first element is not an object")
			}
			structName := "RandomStructItem"
			_, structDefinitions = generateStructFromJSON(firstElem, structName)
			structDefinitions["RandomStruct"] = fmt.Sprintf("type RandomStruct []%s", structName)
		} else {
			log.Fatalf("Error: JSON array is empty")
		}
	} else {
		var jsonObj map[string]interface{}
		err := json.Unmarshal(jsonData, &jsonObj)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON object: %v", err)
		}
		structName := "RandomStruct"
		_, structDefinitions = generateStructFromJSON(jsonObj, structName)
	}

	err := SaveStructToFile(filePath, packageName, structDefinitions)
	if err != nil {
		log.Fatalf("Error saving struct to file: %v", err)
	}
}
