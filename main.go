package main

import (
	"github.com/gin-gonic/gin"
	"wysh-app/controllers"
)

func main() {
	controllers.ConnectDB()

	router := setupRouter()
	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}

}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/home", controllers.GetHome)
	router.GET("/trend", controllers.GetTrend)
	router.GET("/article/detail", controllers.GetArticleDetail)

	return router
}
