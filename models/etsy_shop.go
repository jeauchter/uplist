package models

type EtsyShop struct {
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}

type Result struct {
	ShopID                         int      `json:"shop_id"`
	UserID                         int      `json:"user_id"`
	ShopName                       string   `json:"shop_name"`
	CreateDate                     int      `json:"create_date"`
	CreatedTimestamp               int      `json:"created_timestamp"`
	Title                          string   `json:"title"`
	Announcement                   string   `json:"announcement"`
	CurrencyCode                   string   `json:"currency_code"`
	IsVacation                     bool     `json:"is_vacation"`
	VacationMessage                string   `json:"vacation_message"`
	SaleMessage                    string   `json:"sale_message"`
	DigitalSaleMessage             string   `json:"digital_sale_message"`
	UpdateDate                     int      `json:"update_date"`
	UpdatedTimestamp               int      `json:"updated_timestamp"`
	ListingActiveCount             int      `json:"listing_active_count"`
	DigitalListingCount            int      `json:"digital_listing_count"`
	LoginName                      string   `json:"login_name"`
	AcceptsCustomRequests          bool     `json:"accepts_custom_requests"`
	PolicyWelcome                  string   `json:"policy_welcome"`
	PolicyPayment                  string   `json:"policy_payment"`
	PolicyShipping                 string   `json:"policy_shipping"`
	PolicyRefunds                  string   `json:"policy_refunds"`
	PolicyAdditional               string   `json:"policy_additional"`
	PolicySellerInfo               string   `json:"policy_seller_info"`
	PolicyUpdateDate               int      `json:"policy_update_date"`
	PolicyHasPrivateReceiptInfo    bool     `json:"policy_has_private_receipt_info"`
	HasUnstructuredPolicies        bool     `json:"has_unstructured_policies"`
	PolicyPrivacy                  string   `json:"policy_privacy"`
	VacationAutoreply              string   `json:"vacation_autoreply"`
	URL                            string   `json:"url"`
	ImageURL760x100                string   `json:"image_url_760x100"`
	NumFavorers                    int      `json:"num_favorers"`
	Languages                      []string `json:"languages"`
	IconURLFullxfull               string   `json:"icon_url_fullxfull"`
	IsUsingStructuredPolicies      bool     `json:"is_using_structured_policies"`
	HasOnboardedStructuredPolicies bool     `json:"has_onboarded_structured_policies"`
	IncludeDisputeFormLink         bool     `json:"include_dispute_form_link"`
	IsDirectCheckoutOnboarded      bool     `json:"is_direct_checkout_onboarded"`
	IsEtsyPaymentsOnboarded        bool     `json:"is_etsy_payments_onboarded"`
	IsCalculatedEligible           bool     `json:"is_calculated_eligible"`
	IsOptedInToBuyerPromise        bool     `json:"is_opted_in_to_buyer_promise"`
	IsShopUSBased                  bool     `json:"is_shop_us_based"`
	TransactionSoldCount           int      `json:"transaction_sold_count"`
	ShippingFromCountryISO         string   `json:"shipping_from_country_iso"`
	ShopLocationCountryISO         string   `json:"shop_location_country_iso"`
	ReviewCount                    int      `json:"review_count"`
	ReviewAverage                  int      `json:"review_average"`
}
