package mongoapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	mongodb "github.com/roeeattias/Book-Store/mongoDB/database"
	mongoschemes "github.com/roeeattias/Book-Store/mongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var credentials mongoschemes.Author
	var author mongoschemes.Author

	// binding the credentials got from the request to the credentials object
	if err := c.BindJSON(&credentials); err != nil {
		return
	}

	// checking if the author exists in the database
	author = getAuthorByName(credentials.Username)
	if author.Username == "" {
		c.Status(http.StatusInternalServerError)
		return
	}
	
	// comparing the password entered by the user to the hashed password in the database
	err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(credentials.Password))
    if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// generating a jwt token
	token, err := generateJWTtoken(credentials.Username, author.ID.Hex())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// setting an authorization cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600 * 24 * 30, "", "", false, true)
	c.Status(http.StatusOK)
}

func SignUp(c *gin.Context) {
	var newAuthor mongoschemes.Author

	// pulling author data from the request to the object
	if err := c.BindJSON(&newAuthor); err != nil {
		return
	}

	// hashing the password (Bcrypt, 12 factor)
	password := newAuthor.Password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
    newAuthor.Password = string(bytes)
	newAuthor.PublishedBooks = []primitive.ObjectID{}
	
	// Creating a new user document
	authorInstance, err := mongodb.AuthorCollection.InsertOne(context.Background(), newAuthor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// generating JWT token
	token, err := generateJWTtoken(newAuthor.Username, authorInstance.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	// setting the authorization cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600 * 24 * 30, "", "", false, true)
	c.Status(http.StatusCreated)
}

func PublishBook(c *gin.Context) {
	var newBook mongoschemes.Book

	// pulling the new book data from the request
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	usernameAny, usernameExists := c.Get("username")
	userIdAny, userIdExists := c.Get("userId")

	if !userIdExists || !usernameExists {
		c.Redirect(http.StatusNetworkAuthenticationRequired, "/login")
        return
	}

	// Assert that the username is a string
	username, ok := usernameAny.(string)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Assert that the user id is a string
	userId, ok := userIdAny.(string)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	userIdObject, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// filling neccessary fields to the new book object
	newBook.PublishDate = time.Now()
	newBook.Publisher = username
	newBook.PublisherId = userIdObject
	newBook.Rating = 0

	bookInstance, err := mongodb.BooksCollection.InsertOne(context.Background(), newBook)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	

	// Update the authors collection to add the new book to the author's books array
	update := bson.M{
		"$push": bson.M{
			"published_books": bookInstance.InsertedID,
		},
	}

	filter := bson.M{"_id": userIdObject}

	_, updateErr := mongodb.AuthorCollection.UpdateOne(context.Background(), filter, update)
	if updateErr != nil {
		fmt.Println(updateErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author's books"})
		return
	}

	c.Status(http.StatusCreated)
}