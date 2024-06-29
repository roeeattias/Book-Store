package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()


	fmt.Println("Start listening on port 8080")
	router.Run()
}