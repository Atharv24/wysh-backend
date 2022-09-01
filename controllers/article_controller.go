package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetArticleDetail(c *gin.Context) {
	SampleArticle.ID, _ = strconv.Atoi(c.Request.URL.Query().Get("id"))
	c.JSON(http.StatusOK, SampleArticle)
}
