package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wysh-app/models"
)

func GetArticleDetail(c *gin.Context) {
	//articleId, _ := strconv.Atoi(c.Request.URL.Query().Get("id"))

	article := models.ArticleMini{
		ID:           0,
		Name:         "",
		Brand:        "",
		CurrentPrice: 0,
		BasePrice:    0,
		ImageUrl:     "",
	}
	c.JSON(http.StatusOK, article)
}
