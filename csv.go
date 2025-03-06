package csv

func Reader(filePath string, result any) error {
	return readCsv(filePath, result)
}

func Writer(destination string, fileName string, data any) error {
	return writeCsv(destination, fileName, data)
}
