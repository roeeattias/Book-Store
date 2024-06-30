package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	mongodb "github.com/roeeattias/Book-Store/mongoDB/database"
	mongoapi "github.com/roeeattias/Book-Store/mongoDB/handler"
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
	
	// user authentication endpoints
	router.POST("/login", mongoapi.Login)
	router.POST("/signup", mongoapi.SignUp)

	// book related endpoints
	router.POST("/publishBook", mongoapi.Middleware, mongoapi.PublishBook)
	router.GET("/getBooks", mongoapi.GetBooks)
	router.PATCH("/editBook", mongoapi.Middleware, mongoapi.UpdateBookInformation)
	router.DELETE("/deleteBook", mongoapi.Middleware, mongoapi.DeleteBook)

	fmt.Println("Start listening on post 8080")
	router.Run()
}
