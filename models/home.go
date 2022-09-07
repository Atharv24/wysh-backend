package models

type ArticleMini struct {
	ID           int64  `json:"id"`
	Name         string `json:"title"`
	Brand        string `json:"brand"`
	CurrentPrice int64  `json:"current_price"`
	BasePrice    int64  `json:"base_price"`
	ImageUrl     string `json:"image_url"`
	ArticleUrl   string `json:"article_url"`
}

type TrendMini struct {
	ID       int           `json:"id"`
	Title    string        `json:"title"`
	Articles []ArticleMini `json:"articles"`
}

type HomeResObj struct {
	Trends []TrendMini `json:"trends"`
}
