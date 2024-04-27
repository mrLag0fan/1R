package scanner

import (
	"1R/internal/model"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Data struct {
	Attributes []string
	Data       []model.Entry
}

func NewData(path string) *Data {
	data := Data{}
	err := data.importCSV(path)
	if err != nil {
		return nil
	}
	return &data
}

func (d *Data) importCSV(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Читання CSV-файлу
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	// Перетворення рядків у відповідні типи даних та додавання до даних
	for i, record := range records {
		if i == 0 {
			d.Attributes = record
			continue
		}
		age, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		frequencyPerYear, err := strconv.Atoi(record[1])
		if err != nil {
			return err
		}
		hasPassport := stringToBool(record[3])
		answer := stringToBool(record[4])

		entry := model.Entry{
			Age:               age,
			FrequencyPerYear:  frequencyPerYear,
			NewPlacesAttitude: record[2],
			HasPassport:       hasPassport,
			Answer:            answer,
		}
		d.Data = append(d.Data, entry)
	}

	return nil
}

func ExporyCSV(data Data, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Запис заголовків
	if err := writer.Write([]string{"Age", "FrequencyPerYear", "NewPlacesAttitude", "HasPassport", "Answer"}); err != nil {
		return err
	}

	// Запис даних
	for _, entry := range data.Data {
		hasPassport := ""
		if entry.HasPassport != nil {
			hasPassport = fmt.Sprintf("%t", *entry.HasPassport)
		}
		answer := ""
		if entry.Answer != nil {
			answer = fmt.Sprintf("%t", *entry.Answer)
		}
		if err := writer.Write([]string{
			fmt.Sprintf("%d", entry.Age),
			fmt.Sprintf("%d", entry.FrequencyPerYear),
			entry.NewPlacesAttitude,
			hasPassport,
			answer,
		}); err != nil {
			return err
		}
	}

	return nil
}

func stringToBool(s string) *bool {
	t := true
	f := false
	if s == "Так" {
		return &t
	} else if s == "Ні" {
		return &f
	} else {
		return nil
	}
}
