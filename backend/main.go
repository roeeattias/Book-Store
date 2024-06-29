package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	mongodb "github.com/roeeattias/Book-Store/mongoDB/database"
)

func main() {
	router := gin.Default()
	mongodb.Connect()

	defer func() {
		err := mongodb.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("Start listening on post 8080")
	router.Run()
}
