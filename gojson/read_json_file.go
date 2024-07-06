package gojson

import "os"

func ReadJSONFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}
