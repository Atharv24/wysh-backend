package controllers

import (
	"context"
	"encoding/json"
	"github.com/mindstand/gogm/v2"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"wysh-app/models"
)

// const articleInsertionLimit = 10000
const parallelism = 8
const dbRetries = 5

const hnmProductListingUrlMen = "https://www2.hm.com/en_in/men/shop-by-product/view-all/_jcr_content/main/productlisting_fa5b.display.json" +
	"?sort=stock&image=model&offset=0&page-size=5000"

const hnmProductListingUrlWomen = "https://www2.hm.com/en_in/women/shop-by-product/view-all/_jcr_content/main/productlisting_30ab.display.json" +
	"?sort=stock&image=model&offset=0&page-size=10000"

func getHnmProductAvailabilityUrl(articleStoreID string) string {
	return "https://www2.hm.com/hmwebservices/service/product/in/availability/" + articleStoreID[:7] + ".json"
}

var sizeMaps = struct {
	Upperwear    map[string]*models.Size
	UpperwearMen map[string]*models.Size
}{
	Upperwear:    map[string]*models.Size{},
	UpperwearMen: map[string]*models.Size{},
}

var sizeScales = struct {
	Jeans     []string
	Upperwear []string
}{Jeans: []string{"32", "34", "36", "38", "40", "42", "44", "46"}, Upperwear: []string{"XXS", "XS", "S", "M", "L", "XL", "XXL"}}

var hnmStore = models.Store{
	Name:       "HNM",
	WebsiteUrl: "https://www2.hm.com/en_in/index.html",
}

var hnmBrand = models.Brand{
	Name:            "HNM",
	FirstPartyStore: &hnmStore,
}

func checkVariantEdgeExists(article *models.Article, variation *models.Variation) bool {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)
	query := `MATCH (:Article{id:$articleId}) -[v:HAS_VARIANT]-> (:Variant{id:$variantId}) RETURN v LIMIT 1`
	properties := map[string]interface{}{
		"articleId": article.Id,
		"variantId": variation.Id,
	}
	var variantEdge models.VariantEdge
	err = sess.Query(context.Background(), query, properties, &variantEdge)
	if err != nil {
		return false
	}
	return variantEdge.Id != nil
}

func attachVariationID(variation *models.Variation) {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)
	query := ""
	var properties map[string]interface{}
	if variation.Color != nil && variation.Size != nil {
		query = `MATCH (c:Color{hex_code:$colorCode}), (s:Size{size:$size}), (store:Store{name:"HNM"}) MERGE (c)<-[:HAS_COLOR]-(v:Variation)-[:HAS_SIZE]->(s) MERGE (v)-[:SOLD_AT]->(store) RETURN v`
		properties = map[string]interface{}{
			"colorCode": variation.Color.HexCode,
			"size":      variation.Size.Size,
		}
		err = sess.Query(context.Background(), query, properties, variation)

	} else if variation.Color != nil {
		query = `MATCH (c:Color{hex_code:$colorCode}), (store:Store{name:"HNM"}) MERGE (c)<-[:HAS_COLOR]-(v:Variation)-[:SOLD_AT]->(store) RETURN v`
		properties = map[string]interface{}{
			"colorCode": variation.Color.HexCode,
		}
		err = sess.Query(context.Background(), query, properties, variation)

	} else {
		query = `MATCH (store:Store{name:"HNM"}) MERGE (v:Variation)-[:SOLD_AT]->(store) RETURN v`
		properties = make(map[string]interface{})
		err = sess.Query(context.Background(), query, properties, variation)
	}

	if err != nil {
		err = sess.Save(context.Background(), variation)
		if err != nil {
			for i := 0; i < dbRetries; i++ {
				err = sess.Save(context.Background(), variation)
				if err == nil {
					break
				}
			}
			if err != nil {
				panic(err)
			}
		}
	}
}

func getTag(tagStr string) *models.Tag {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)

	var tag models.Tag
	err = sess.Query(context.Background(), "MERGE (t:Tag{name:$tag}) RETURN t", map[string]interface{}{
		"tag": tagStr,
	}, &tag)

	if err != nil {
		tag.Name = tagStr
		err = sess.Save(context.Background(), &tag)
		if err != nil {
			for i := 0; i < dbRetries; i++ {
				err = sess.Save(context.Background(), &tag)
				if err == nil {
					break
				}
			}
			if err != nil {
				panic(err)
			}
		}
	}

	return &tag
}

func getSize(sizeStr string) *models.Size {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)

	var size models.Size
	err = sess.Query(context.Background(), "MERGE (s:Size{size:$size}) RETURN s LIMIT 1", map[string]interface{}{
		"size": sizeStr,
	}, &size)

	if err != nil {
		size.Size = sizeStr
		err = sess.Save(context.Background(), &size)
		if err != nil {
			panic(err)
		}
	}

	return &size
}

func getColor(colorCode string) *models.Color {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)

	var color models.Color
	err = sess.Query(context.Background(), "MERGE (c:Color{hex_code:$colorCode}) RETURN c LIMIT 1", map[string]interface{}{
		"colorCode": colorCode,
	}, &color)

	if err != nil {
		color.HexCode = colorCode
		err = sess.Save(context.Background(), &color)
		if err != nil {
			for i := 0; i < dbRetries; i++ {
				err = sess.Save(context.Background(), &color)
				if err == nil {
					break
				}
			}
		}
		if err != nil {
			panic(err)
		}
	}
	return &color
}

func getStoreBrand() {
	sess, err := gogm.G().NewSessionV2(gogm.SessionConfig{AccessMode: gogm.AccessModeWrite})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func(sess gogm.SessionV2) {
		err = sess.Close()
		if err != nil {
			panic(err)
		}
	}(sess)
	err = sess.Save(context.Background(), &hnmStore)
	if err != nil {
		return
	}

	err = sess.Save(context.Background(), &hnmBrand)
	if err != nil {
		return
	}
}

func setSizeMap() {
	for i := 0; i < 7; i++ {
		j := i + 8
		sizeCode := strconv.Itoa(j)
		_sizeCode := strconv.Itoa(j + 1)

		for len(sizeCode) < 3 {
			sizeCode = "0" + sizeCode
		}
		for len(_sizeCode) < 3 {
			_sizeCode = "0" + _sizeCode
		}
		size := getSize(sizeScales.Upperwear[i])

		sizeMaps.Upperwear[sizeCode] = size
		sizeMaps.UpperwearMen[_sizeCode] = size
	}
}

func sendRequest(url string) []byte {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("User-Agent", "web-browser")
	res, _ := http.DefaultClient.Do(req)
	bodyBytes, _ := io.ReadAll(res.Body)
	return bodyBytes
}

func PullHnmData() {
	hnmStore.FirstPartyBrand = &hnmBrand
	getStoreBrand()
	setSizeMap()

	var resObjWomen models.HnmAllProducts
	bodyBytes := sendRequest(hnmProductListingUrlWomen)
	_ = json.Unmarshal(bodyBytes, &resObjWomen)

	var resObjMen models.HnmAllProducts
	bodyBytes = sendRequest(hnmProductListingUrlMen)
	_ = json.Unmarshal(bodyBytes, &resObjMen)

	products := append(resObjWomen.Products, resObjMen.Products...)
	chunkSize := len(products) / parallelism
	for i := 1; i < parallelism; i++ {
		go updateArticles(products[i*chunkSize : (i+1)*chunkSize])
	}
	updateArticles(products[(parallelism-1)*chunkSize:])
}

func updateArticles(products []models.HnmProduct) {
	for _, product := range products {
		updateOrCreateArticle(product)
	}
}

func updateOrCreateArticle(product models.HnmProduct) {
	article := getArticleByStoreArticleCode(product.ArticleCode[:7], hnmStore.Name)
	if article == nil {
		imageUrl := product.Image[0].Src
		productUrl := product.Link
		article = &models.Article{
			Name:              product.Title,
			Brand:             &hnmBrand,
			DefaultImageUrl:   encodeUrl("https:" + imageUrl),
			DefaultProductUrl: encodeUrl("https://www2.hm.com" + productUrl),
		}
		var tags []string
		for _, tag := range strings.Split(product.Category, "_") {
			tags = append(tags, tag)
		}
		for _, tag := range strings.Split(product.Title, " ") {
			tags = append(tags, tag)
		}
		var _tags []*models.Tag
		tagSet := make(map[string]bool)
		re, err := regexp.Compile(`\W`)
		if err != nil {
			panic(err)
		}

		for _, tag := range tags {
			tagStr := re.ReplaceAllString(tag, "")
			tagLower := strings.ToLower(tagStr)
			if !tagSet[tagLower] {
				tagSet[tagLower] = true
				_tag := getTag(tagLower)
				_tags = append(_tags, _tag)
			}
		}
		article.Tags = _tags
	}
	attachVariations(&product, article)
	insertArticle(article)
}

func attachVariations(product *models.HnmProduct, article *models.Article) {
	resObj := struct {
		Availability  []string `json:"availability"`
		FewPiecesLeft []string `json:"fewPieceLeft"`
	}{}
	bodyBytes := sendRequest(getHnmProductAvailabilityUrl(product.ArticleCode))
	_ = json.Unmarshal(bodyBytes, &resObj)

	colorSizeMap := map[string][]string{}
	for _, availability := range resObj.Availability {
		articleColorCode := availability[:10]
		if colorSizeMap[articleColorCode] == nil {
			colorSizeMap[articleColorCode] = []string{}
		}
		colorSizeMap[articleColorCode] = append(colorSizeMap[articleColorCode], availability[10:])
	}

	_tags := strings.Split(product.Category, "_")
	targetGender := _tags[0]
	category := ""
	if len(_tags) > 1 {
		category = _tags[1]
	}
	var subcategories []string
	if len(_tags) > 2 {
		subcategories = _tags[2:]
	}

	_sizeMap := getSizeMap(targetGender, category, subcategories)

	priceRegExp, _ := regexp.Compile("[^0-9]")
	price, _ := strconv.Atoi(priceRegExp.ReplaceAllString(product.Price, ""))

	_colors := product.Swatches

	var edges []*models.VariantEdge
	variation := models.Variation{
		Store: &hnmStore,
	}
	if len(_colors) == 0 {
		attachVariationID(&variation)
		edgeExists := checkVariantEdgeExists(article, &variation)
		if !edgeExists {
			edge := models.VariantEdge{
				Start:          article,
				End:            &variation,
				Available:      len(resObj.Availability) > 0,
				ColorName:      "",
				CurrentPrice:   price,
				BasePrice:      price,
				VariantStoreID: product.ArticleCode,
				VariantUrl:     encodeUrl("https://www2.hm.com" + product.Link),
			}
			edges = append(edges, &edge)
		}
		return
	}
	validSizeSystem := _sizeMap != nil

	for _, color := range _colors {
		if !validSizeSystem {
			break
		}
		articleColorCode := strings.Split(color.ArticleLink, ".")[1]
		sizes := colorSizeMap[articleColorCode]
		for _, sizeCode := range sizes {
			size := _sizeMap[sizeCode]
			if size == nil {
				validSizeSystem = false
				break
			}
		}
	}
	for _, color := range _colors {
		articleColorCode := strings.Split(color.ArticleLink, ".")[1]
		clr := getColor(color.ColorCode)
		if !validSizeSystem {
			variation = models.Variation{Store: &hnmStore, Color: clr}
			attachVariationID(&variation)
			sizes := colorSizeMap[articleColorCode]
			edgeExists := checkVariantEdgeExists(article, &variation)
			if !edgeExists {
				edge := models.VariantEdge{
					Start:          article,
					End:            &variation,
					Available:      len(sizes) > 0,
					ColorName:      color.ColorName,
					CurrentPrice:   price,
					BasePrice:      price,
					VariantStoreID: articleColorCode,
					VariantUrl:     encodeUrl("https://www2.hm.com" + color.ArticleLink),
				}
				edges = append(edges, &edge)
			}
		} else {
			sizes := colorSizeMap[articleColorCode]
			if sizes != nil && len(sizes) > 0 {
				for _, sizeCode := range sizes {
					variation = models.Variation{Store: &hnmStore, Color: clr, Size: _sizeMap[sizeCode]}
					attachVariationID(&variation)
					edgeExists := checkVariantEdgeExists(article, &variation)
					if !edgeExists {
						edge := models.VariantEdge{
							Start:          article,
							End:            &variation,
							Available:      true,
							CurrentPrice:   price,
							BasePrice:      price,
							VariantStoreID: articleColorCode,
							VariantUrl:     encodeUrl("https://www2.hm.com" + color.ArticleLink),
						}
						edges = append(edges, &edge)
					}
				}
			}
		}
	}

	article.Variations = edges
}

func getSizeMap(gender string, category string, subcategories []string) map[string]*models.Size {
	if gender == "ladies" || gender == "men" {
		if category == "tops" || category == "blazers" || category == "cardigansjumpers" ||
			category == "dresses" || category == "hoodiesswetshirts" || category == "jacketscoats" ||
			category == "jumpsuits" || category == "knitwear" || category == "licence" ||
			category == "shirtsblouses" || (gender == "ladies" && category == "trousers") ||
			category == "tshirtstanks" || category == "shirts" ||
			category == "shirt" || category == "blazersuits" ||
			(category == "basics" && len(subcategories) > 0 && subcategories[0] != "lingerie") ||
			(category == "premium" && len(subcategories) > 0 && subcategories[0] == "tops") ||
			(category == "skirts" && len(subcategories) > 0 && (subcategories[0] == "highwaisted" ||
				subcategories[0] == "shortskirts" || subcategories[0] == "denim")) {
			if gender == "men" {
				return sizeMaps.UpperwearMen
			}
			return sizeMaps.Upperwear
		}
	}
	return nil
}

func encodeUrl(s string) string {
	sepStr := strings.Split(s, "?")
	if len(sepStr) <= 1 {
		return s
	}
	return sepStr[0] + "?" + url.QueryEscape(strings.Join(sepStr[1:], "?"))
}
