package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/primadi/webapi/utils"
)

// SET using -ldflags="-X 'main.BuildTime=$(date)'"
var BuildTime string

func SetupRouter() *gin.Engine {
	router := utils.NewRouter(BuildTime)

	return router
}

func main() {
	godotenv.Load()
	router := SetupRouter()

	utils.ListenAndServe(router)
}
