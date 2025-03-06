package csv

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

func writeCsv(destination string, fileName string, input any) error {
	inputValue := reflect.ValueOf(input)
	if inputValue.Kind() != reflect.Ptr || inputValue.Kind() != reflect.Slice {
		return errors.New("input must be a pointer to the slice")
	}

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		return err
	}

	filePath := filepath.Join(destination, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if inputValue.Len() == 0 {
		return nil
	}

	firstElement := inputValue.Index(0)
	elementType := firstElement.Type()

	var headers []string
	for i := 0; i < elementType.NumField(); i++ {
		field := elementType.Field(i)
		headers = append(headers, field.Name)
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	for i := 0; i < inputValue.Len(); i++ {
		item := inputValue.Index(i)
		var row []string

		for j := 0; j < len(headers); j++ {
			field := item.FieldByName(headers[j])
			row = append(row, field.String())
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
