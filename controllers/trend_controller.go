package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wysh-app/models"
)

func GetTrend(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.URL.Query().Get("id"))
	trend := models.TrendMini{ID: id, Articles: []models.ArticleMini{}}
	SampleArticle.ID = 1
	for i := 0; i < 20; i++ {
		trend.Articles = append(trend.Articles, SampleArticle)
		SampleArticle.ID++
	}
	c.JSON(http.StatusOK, trend)
}
