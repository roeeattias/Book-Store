package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	mongodb "github.com/roeeattias/Book-Store/mongoDB/database"
	mongoapi "github.com/roeeattias/Book-Store/mongoDB/handler"
)

func main() {
	router := gin.Default()
	mongodb.Connect()
	
	// Define your CORS configuration
    corsConfig := cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // Allowed origin(s)
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Allowed methods
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Allowed headers
        ExposeHeaders:    []string{"Content-Length"}, // Expose specific headers to the client
        AllowCredentials: true, // Allow credentials (cookies, authorization headers)
    }

    // Apply the CORS middleware to the router
    router.Use(cors.New(corsConfig))

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
	router.POST("/buyBook", mongoapi.BuyBook)
	router.POST("/getAuthors", mongoapi.GetAuthors)

	fmt.Println("Start listening on post 8080")
	router.Run()
}
