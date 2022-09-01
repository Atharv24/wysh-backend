package controllers

import (
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"wysh-app/models"
)

const hnmProductListingUrl = "https://www2.hm.com/en_in/men/shop-by-product/view-all/_jcr_content/main/productlisting_fa5b.display.json" +
	"?sort=stock&image=model&offset=0&page-size=5000"

func getHnmProductAvailablityUrl(articleStoreID string) string {
	return "https://www2.hm.com/hmwebservices/service/product/in/availability/" + articleStoreID[:7] + ".json"
}

func sendRequest(url string) []byte {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("User-Agent", "web-browser")
	res, _ := http.DefaultClient.Do(req)
	bodyBytes, _ := io.ReadAll(res.Body)
	return bodyBytes
}

func PullHnmData() {
	var resObj models.HnmAllProducts
	bodyBytes := sendRequest(hnmProductListingUrl)
	_ = json.Unmarshal(bodyBytes, &resObj)
	var articles []models.Article

	store := models.Store{
		Title:      "HNM",
		WebsiteUrl: "https://www2.hm.com/en_in/index.html",
	}
	brand := models.Brand{
		Title:           "HNM",
		FirstPartyStore: store,
	}

	for _, product := range resObj.Products {
		var tags []models.Tag
		for _, tagStr := range strings.Split(product.Category, "_") {
			tags = append(tags, models.Tag{Title: strings.ToLower(tagStr)})
		}

		//variations := getVariations(product)

		imageUrl := product.Image[0].Src
		productUrl := product.Link
		article := models.Article{
			Title:             product.Title,
			ArticleStoreID:    product.ArticleCode,
			Brand:             brand,
			Tags:              tags,
			DefaultImageUrl:   encodeUrl(imageUrl),
			DefaultProductUrl: encodeUrl(productUrl),
			//Variations:        variations,
		}
		articles = append(articles, article)
	}
	fmt.Println(articles)
}

func getVariations(product models.HnmProduct) []models.Variation {
	resObj := struct {
		Availability []string `json:"availability"`
	}{}
	bodyBytes := sendRequest(getHnmProductAvailablityUrl(product.ArticleCode))
	_ = json.Unmarshal(bodyBytes, &resObj)
	tags := strings.Split(product.Category, "_")
	targetGender := tags[0]
	categoru := tags[1]

	return nil
}

func encodeUrl(s string) string {
	sepStr := strings.Split(s, "?")
	return sepStr[0] + "?" + url.QueryEscape(strings.Join(sepStr[1:], ""))
}
