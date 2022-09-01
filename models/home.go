package models

import "strings"

type ArticleMini struct {
	ID           int    `json:"id"'`
	Title        string `json:"title"`
	Brand        string `json:"brand"`
	CurrentPrice int    `json:"current_price"`
	BasePrice    int    `json:"base_price"`
	ImageLinks   string
	ImageUrl     string `json:"image_url"`
}

func (a *ArticleMini) parseImageUrl() {
	a.ImageUrl = "https://" + strings.Split(a.ImageLinks, ",")[0]
}

type TrendMini struct {
	ID       int           `json:"id"`
	Title    string        `json:"title"`
	Articles []ArticleMini `json:"articles"`
}

func (t *TrendMini) ParseImages() {
	for i := range t.Articles {
		t.Articles[i].parseImageUrl()
	}
}

type HomeResObj struct {
	Trends []TrendMini `json:"trends"`
}
