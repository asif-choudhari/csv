package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
)

func readCsv(filePath string, result any) error {
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr || resultValue.Elem().Kind() != reflect.Slice {
		return errors.New("result must be a pointer to the slice")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := csv.NewReader(file)

	headers, err := encoder.Read()
	if err != nil {
		return err
	}

	records, err := encoder.ReadAll()
	if err != nil {
		return err
	}

	for index, record := range records {
		if len(record) != len(headers) {
			return fmt.Errorf("invalid csv format at line %d", index+2)
		}
	}

	sliceType := resultValue.Elem().Type().Elem()

	slice := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, len(records))

	for _, record := range records {
		item := reflect.New(sliceType).Elem()

		for i, header := range headers {
			field := item.FieldByName(header)
			if field.IsValid() && field.CanSet() {
				field.SetString(record[i])
			}
		}

		slice = reflect.Append(slice, item)
	}

	resultValue.Elem().Set(slice)
	return nil
}
