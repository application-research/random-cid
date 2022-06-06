package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/application-research/random-cid/cli"
	"github.com/gin-gonic/gin"
)

func getCidV1(ctx *gin.Context) {
	version := 1
	c, err := cli.NewCid(version)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	ctx.String(http.StatusOK, c.String()+"\n")
}

func getCidV0(ctx *gin.Context) {
	version := 0
	c, err := cli.NewCid(version)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	ctx.String(http.StatusOK, c.String()+"\n")
}

func main() {
	router := gin.Default()
	router.GET("/", getCidV1)
	router.GET("/v1", getCidV1)
	router.GET("/v0", getCidV0)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
