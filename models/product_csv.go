package models

type ProductCSV struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Option1Name         string   `json:"option_1_name"`
	Option1Value        string   `json:"option_1_value"`
	Option2Name         string   `json:"option_2_name"`
	Option2Value        string   `json:"option_2_value"`
	Option3Name         string   `json:"option_3_name"`
	Option3Value        string   `json:"option_3_value"`
	SKU                 string   `json:"sku"`
	GTIN                string   `json:"gtin"`
	ASIN                string   `json:"asin"`
	Quantity            string   `json:"quantity"`
	Price               float64  `json:"price"`
	ImageLinks          string   `json:"image_link"`
	AdditionalImageLink []string `json:"additional_image_link"`
	Brand               string   `json:"brand"`
	Tags                []string `json:"tags"`
	Category            string   `json:"category"`
	Weight              string   `json:"weight"`
	WeightUnit          string   `json:"weight_unit"`
	Height              string   `json:"height"`
	Width               string   `json:"width"`
	Length              string   `json:"length"`
	DimensionsUnits     string   `json:"dimensions_units"`
}
