package repository

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/jeremyauchter/uplist/models"
)

func GetResources() []models.Resource {
	// Fetch resources from database and return
	if len(os.Args) < 2 {
		log.Fatal("Please provide a CSV file as a command line argument")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	csvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var record models.Resource
	var records []models.Resource

	for _, each := range csvData {
		record.ID = each[0]
		record.Title = each[1]
		record.Description = each[2]
		record.Option1Name = each[3]
		record.Option1Value = each[4]
		record.Option2Name = each[5]
		record.Option2Value = each[6]
		record.Option3Name = each[7]
		record.Option3Value = each[8]
		record.SKU = each[9]
		record.GTIN = each[10]
		record.ASIN = each[11]
		record.Quantity = each[12]
		record.Price = each[13]
		record.ImageLink = each[14]
		record.AdditionalImageLink = each[15]
		record.Brand = each[16]
		record.Tags = each[17]
		record.Category = each[18]
		record.Weight = each[19]
		record.WeightUnit = each[20]
		record.Height = each[21]
		record.Width = each[22]
		record.Length = each[23]
		record.DimensionsUnits = each[24]
		// Assign the rest of the fields
		records = append(records, record)
	}

	return records
}
