package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wysh-app/models"
)

func GetHome(c *gin.Context) {
	var trends []models.TrendMini
	article := getArticleDetailByVariationId(7929)

	for i := 0; i < 5; i++ {
		trend := models.TrendMini{
			ID:       i,
			Title:    "TREND",
			Articles: []models.ArticleMini{},
		}
		for j := 0; j < 5; j++ {
			trend.Articles = append(trend.Articles, *article)
			article.ID++
		}
		trends = append(trends, trend)
	}

	obj := models.HomeResObj{Trends: trends}
	c.JSON(http.StatusOK, obj)
}
