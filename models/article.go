package models

type Article struct {
	ID                int    `json:"id"`
	Title             string `json:"title,omitempty"`
	Description       string
	Brand             `json:"brand"`
	Tags              []Tag
	DefaultProductUrl string
	DefaultImageUrl   string
	Variations        []Variation
	ArticleStoreID    string
}

type Tag struct {
	ID    int
	Title string
}

type Color struct {
	ID      int
	Title   string
	HexCode string
}

type Pattern struct {
	ID     int
	Title  string
	Colors []Color
}

type Size struct {
	ID   int
	Size Sizes
}

type Sizes string

const (
	XS  Sizes = "XS"
	S   Sizes = "S"
	M   Sizes = "M"
	L   Sizes = "L"
	XL  Sizes = "XL"
	XXL Sizes = "XXL"
)

type Variation struct {
	ID           int
	CurrentPrice int
	BasePrice    int
	Pattern
	Size
	Store
}

type Brand struct {
	ID              int
	Title           string
	FirstPartyStore Store
}

type Store struct {
	ID         int
	Title      string
	WebsiteUrl string
}
