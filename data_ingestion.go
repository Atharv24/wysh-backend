package main

import "wysh-app/controllers"

func main() {
	controllers.ConnectDB()
	controllers.PullHnmData()
}
