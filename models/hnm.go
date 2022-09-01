package models

type HnmAllProducts struct {
	Products []HnmProduct `json:"products"`
}

type HnmProduct struct {
	ArticleCode string `json:"articleCode"`
	OnClick     string `json:"onClick"`
	Link        string `json:"link"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Image       []struct {
		Src          string `json:"src"`
		DataAltImage string `json:"dataAltImage"`
		Alt          string `json:"alt"`
		DataAltText  string `json:"dataAltText"`
	} `json:"image"`
	LegalText                 string `json:"legalText"`
	PromotionalMarkerText     string `json:"promotionalMarkerText"`
	ShowPromotionalClubMarker bool   `json:"showPromotionalClubMarker"`
	ShowPriceMarker           bool   `json:"showPriceMarker"`
	FavouritesTracking        string `json:"favouritesTracking"`
	FavouritesSavedText       string `json:"favouritesSavedText"`
	FavouritesNotSavedText    string `json:"favouritesNotSavedText"`
	MarketingMarkerText       string `json:"marketingMarkerText"`
	MarketingMarkerType       string `json:"marketingMarkerType"`
	MarketingMarkerCss        string `json:"marketingMarkerCss"`
	Price                     string `json:"price"`
	RedPrice                  string `json:"redPrice"`
	YellowPrice               string `json:"yellowPrice"`
	BluePrice                 string `json:"bluePrice"`
	ClubPriceText             string `json:"clubPriceText"`
	SellingAttribute          string `json:"sellingAttribute"`
	SwatchesTotal             string `json:"swatchesTotal"`
	Swatches                  []struct {
		ColorCode   string `json:"colorCode"`
		ArticleLink string `json:"articleLink"`
		ColorName   string `json:"colorName"`
	} `json:"swatches"`
	PreAccessStartDate string        `json:"preAccessStartDate"`
	PreAccessEndDate   string        `json:"preAccessEndDate"`
	PreAccessGroups    []interface{} `json:"preAccessGroups"`
	OutOfStockText     string        `json:"outOfStockText"`
	ComingSoon         string        `json:"comingSoon"`
	BrandName          string        `json:"brandName"`
	DamStyleWith       string        `json:"damStyleWith"`
}
