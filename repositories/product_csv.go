package repositories

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jeremyauchter/uplist/models"
)

type ProductCSVRepository interface {
	GetProducts(csvPath string) []models.ProductCSV
}

type ProductCSVRepo struct{}

func NewProductCSVRepository() ProductCSVRepository {
	return &ProductCSVRepo{}
}

func (r *ProductCSVRepo) GetProducts(csvPath string) []models.ProductCSV {

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.Read()

	csvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var record models.ProductCSV
	var records []models.ProductCSV

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
		quantity, err := strconv.Atoi(each[12])
		if err != nil {
			log.Fatal(err)
		}
		record.Quantity = quantity
		price, err := strconv.ParseFloat(strings.Replace(each[13], ",", ".", -1), 64)
		if err != nil {
			log.Fatal(err)
		}
		record.Price = price
		var imageLinks []string
		record.ImageLinks = each[14]
		record.AdditionalImageLink = append(imageLinks, strings.Split(each[15], ",")...)
		record.Brand = each[16]
		var tags []string
		record.Tags = append(tags, strings.Split(each[17], ",")...)
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
