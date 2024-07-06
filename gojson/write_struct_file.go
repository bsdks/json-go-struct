package gojson

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

func SaveStructToFile(filename, packageName string, structDefinitions map[string]string) error {
	err := os.MkdirAll(packageName, os.ModePerm)
	if err != nil {
		return err
	}

	fileContent := fmt.Sprintf("package %s\n\n", packageName)
	for _, structDefinition := range structDefinitions {
		fileContent += structDefinition + "\n\n"
	}

	formattedContent, err := format.Source([]byte(fileContent))
	if err != nil {
		return err
	}

	dtoFilePath := filepath.Join(packageName, filename)
	err = os.WriteFile(dtoFilePath, formattedContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
