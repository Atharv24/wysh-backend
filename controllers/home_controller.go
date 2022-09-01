package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wysh-app/models"
)

func GetHome(c *gin.Context) {
	var trends []models.TrendMini

	for i := 0; i < 5; i++ {
		trend := models.TrendMini{
			ID:       i,
			Title:    "TREND",
			Articles: []models.ArticleMini{},
		}
		for j := 0; j < 5; j++ {
			trend.Articles = append(trend.Articles, SampleArticle)
			SampleArticle.ID++
		}
		trends = append(trends, trend)
	}

	obj := models.HomeResObj{Trends: trends}
	c.JSON(http.StatusOK, obj)
}

var SampleArticle = models.ArticleMini{
	ID:           1,
	Title:        "AABC",
	Brand:        "AS",
	CurrentPrice: 1000,
	BasePrice:    1100,
	ImageLinks:   "",
	ImageUrl:     "",
}
