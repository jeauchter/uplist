package models

type WhoMade string

type WhenMade string

const (
	IDid          WhoMade  = "i_did"
	SomeoneElse   WhoMade  = "someone_else"
	Collective    WhoMade  = "collective"
	MadeToOrder   WhenMade = "made_to_order"
	Year2020_2024 WhenMade = "2020_2024"
	Year2010_2019 WhenMade = "2010_2019"
	Year2005_2009 WhenMade = "2005_2009"
	Before2005    WhenMade = "before_2005"
	Year2000_2004 WhenMade = "2000_2004"
	Decade1990s   WhenMade = "1990s"
	Decade1980s   WhenMade = "1980s"
	Decade1970s   WhenMade = "1970s"
	Decade1960s   WhenMade = "1960s"
	Decade1950s   WhenMade = "1950s"
	Decade1940s   WhenMade = "1940s"
	Decade1930s   WhenMade = "1930s"
	Decade1920s   WhenMade = "1920s"
	Decade1910s   WhenMade = "1910s"
	Decade1900s   WhenMade = "1900s"
	Decade1800s   WhenMade = "1800s"
	Decade1700s   WhenMade = "1700s"
	Before1700    WhenMade = "before_1700"
)

type EtsyListingRequest struct {
	ShopID      string   `json:"shop_id"`
	Quantity    string   `json:"quantity"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	WhoMade     WhoMade  `json:"who_made"`
	WhenMade    WhenMade `json:"when_made"`
	TaxonomyID  int      `json:"taxonomy_id"`
	Tags        string   `json:"tags"`
	ImageIDs    string   `json:"image_ids"`
}

type EtsyListingResponse struct {
	ListingID                   int          `json:"listing_id"`
	UserID                      int          `json:"user_id"`
	ShopID                      int          `json:"shop_id"`
	Title                       string       `json:"title"`
	Description                 string       `json:"description"`
	State                       string       `json:"state"`
	CreationTimestamp           int          `json:"creation_timestamp"`
	CreatedTimestamp            int          `json:"created_timestamp"`
	EndingTimestamp             int          `json:"ending_timestamp"`
	OriginalCreationTimestamp   int          `json:"original_creation_timestamp"`
	LastModifiedTimestamp       int          `json:"last_modified_timestamp"`
	UpdatedTimestamp            int          `json:"updated_timestamp"`
	StateTimestamp              int          `json:"state_timestamp"`
	Quantity                    int          `json:"quantity"`
	ShopSectionID               int          `json:"shop_section_id"`
	FeaturedRank                int          `json:"featured_rank"`
	URL                         string       `json:"url"`
	NumFavorers                 int          `json:"num_favorers"`
	NonTaxable                  bool         `json:"non_taxable"`
	IsTaxable                   bool         `json:"is_taxable"`
	IsCustomizable              bool         `json:"is_customizable"`
	IsPersonalizable            bool         `json:"is_personalizable"`
	PersonalizationIsRequired   bool         `json:"personalization_is_required"`
	PersonalizationCharCountMax int          `json:"personalization_char_count_max"`
	PersonalizationInstructions string       `json:"personalization_instructions"`
	ListingType                 string       `json:"listing_type"`
	Tags                        []string     `json:"tags"`
	Materials                   []string     `json:"materials"`
	ShippingProfileID           int          `json:"shipping_profile_id"`
	ReturnPolicyID              int          `json:"return_policy_id"`
	ProcessingMin               int          `json:"processing_min"`
	ProcessingMax               int          `json:"processing_max"`
	WhoMade                     string       `json:"who_made"`
	WhenMade                    string       `json:"when_made"`
	IsSupply                    bool         `json:"is_supply"`
	ItemWeight                  int          `json:"item_weight"`
	ItemWeightUnit              string       `json:"item_weight_unit"`
	ItemLength                  int          `json:"item_length"`
	ItemWidth                   int          `json:"item_width"`
	ItemHeight                  int          `json:"item_height"`
	ItemDimensionsUnit          string       `json:"item_dimensions_unit"`
	IsPrivate                   bool         `json:"is_private"`
	Style                       []string     `json:"style"`
	FileData                    string       `json:"file_data"`
	HasVariations               bool         `json:"has_variations"`
	ShouldAutoRenew             bool         `json:"should_auto_renew"`
	Language                    string       `json:"language"`
	Price                       ListingPrice `json:"price"`
	TaxonomyID                  int          `json:"taxonomy_id"`
}

type ListingPrice struct {
	Amount       int    `json:"amount"`
	Divisor      int    `json:"divisor"`
	CurrencyCode string `json:"currency_code"`
}

type EtsyListingInventoryRequest struct {
	Products           []EtsyProduct `json:"products"`
	PriceOnProperty    []int         `json:"price_on_property"`
	QuantityOnProperty []int         `json:"quantity_on_property"`
	SkuOnProperty      []int         `json:"sku_on_property"`
}

type EtsyProduct struct {
	SKU            string              `json:"sku"`
	PropertyValues []EtsyPropertyValue `json:"property_values"`
	Offerings      []EtsyOffering      `json:"offerings"`
}

type EtsyPropertyValue struct {
	PropertyID   int      `json:"property_id"`
	ValueIDs     []int    `json:"value_ids"`
	ScaleID      int      `json:"scale_id"`
	PropertyName string   `json:"property_name"`
	Values       []string `json:"values"`
}

type EtsyOffering struct {
	Price     int  `json:"price"`
	Quantity  int  `json:"quantity"`
	IsEnabled bool `json:"is_enabled"`
}

type EtsyListingImageRequest struct {
	Image          string `json:"image"`
	ListingImageID int    `json:"listing_image_id"`
	Overwrite      bool   `json:"overwrite" default:"false"`
	IsWatermarked  bool   `json:"is_watermarked" default:"false"`
	AltText        string `json:"alt_text" default:""`
}

type EtsyListingImageResponse struct {
	ListingID         int    `json:"listing_id"`
	ListingImageID    int    `json:"listing_image_id"`
	HexCode           string `json:"hex_code"`
	Red               int    `json:"red"`
	Green             int    `json:"green"`
	Blue              int    `json:"blue"`
	Hue               int    `json:"hue"`
	Saturation        int    `json:"saturation"`
	Brightness        int    `json:"brightness"`
	IsBlackAndWhite   bool   `json:"is_black_and_white"`
	CreationTimestamp int    `json:"creation_tsz"`
	CreatedTimestamp  int    `json:"created_timestamp"`
	Rank              int    `json:"rank"`
	URL75x75          string `json:"url_75x75"`
	URL170x135        string `json:"url_170x135"`
	URL570xN          string `json:"url_570xN"`
	URLFullXFull      string `json:"url_fullxfull"`
	FullHeight        int    `json:"full_height"`
	FullWidth         int    `json:"full_width"`
	AltText           string `json:"alt_text"`
}
