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
	
	// converting the user id string into object id
	userIdObject, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// filling neccessary fields to the new book object
	newBook.ID = primitive.ObjectID{}
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
	
	newBookId, ok := bookInstance.InsertedID.(primitive.ObjectID)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	newBook.ID = newBookId
	c.JSON(http.StatusCreated, newBook)
}

func GetBooks(c *gin.Context) {
	var books []mongoschemes.Book
	
	// Find documents using the filter
	filter := bson.D{{}}
	cursor, err := mongodb.BooksCollection.Find(context.Background(), filter)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and decode each document into the books slice
	for cursor.Next(context.Background()) {
		var book mongoschemes.Book
		if err := cursor.Decode(&book); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// Check for errors after iterating through the cursor
	if err := cursor.Err(); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Send the response
	c.JSON(http.StatusOK, books)
}

func UpdateBookInformation(c *gin.Context) {
	var updatedBookData mongoschemes.Book
	
	// pulling the new updated book data from the request
	if err := c.BindJSON(&updatedBookData); err != nil {
		return
	}

	// getting the currently authenticated user from the moddlewere
	userIdAny, userIdExists := c.Get("userId")
	if !userIdExists {
		c.Redirect(http.StatusNetworkAuthenticationRequired, "/login")
        return
	}

	// Assert that the user id is a string
	userId, ok := userIdAny.(string)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	
	// converting the user id string into object id
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// getting the book information from the database bafore updating it
	// to ensure authenticity of the data
	currentBook, err := getBookById(updatedBookData.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// checking if the authenticated user is indeed the author that published the book
	if userObjectId != currentBook.PublisherId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not your book to update"})
		return
	}

	// Update the book instance with the new data
	update := bson.M{
		"$set": bson.M{
			"title": updatedBookData.Title,
			"author": updatedBookData.Author,
			"quantity": updatedBookData.Quantity,
		},
	}

	// updating the book data
	filter := bson.M{"_id": currentBook.ID}
	_, updateErr := mongodb.BooksCollection.UpdateOne(context.Background(), filter, update)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the book data"})
		return
	}

	c.Status(http.StatusOK)
}

func DeleteBook(c *gin.Context) {
	var bookToDelete mongoschemes.Book

	// pulling the book object to delete from the request
	if err := c.BindJSON(&bookToDelete); err != nil {
		return
	}

	// getting the currently authenticated user from the moddlewere layer
	userIdAny, userIdExists := c.Get("userId")
	if !userIdExists {
		c.Redirect(http.StatusNetworkAuthenticationRequired, "/login")
        return
	}

	// Assert that the user id is a string
	userId, ok := userIdAny.(string)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	
	// converting the user id string into object id
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	
	// getting the book information from the database bafore updating it
	// to ensure authenticity of the data
	currentBook, err := getBookById(bookToDelete.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// checking if the authenticated user is indeed the author that published the book
	if userObjectId != currentBook.PublisherId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not your book to delete"})
		return
	}

	// Define the filter to match the ObjectId
    filter := bson.M{"_id": bookToDelete.ID}

    // Perform the delete operation
    _, err = mongodb.BooksCollection.DeleteOne(context.TODO(), filter)
    if err != nil {
        c.Status(http.StatusInternalServerError)
		return
    }

	c.Status(http.StatusOK)
}

func BuyBook(c *gin.Context) {
	var book mongoschemes.Book
	
	// pulling the new updated book data from the request
	if err := c.BindJSON(&book); err != nil {
		return
	}

	// getting the book information from the database bafore updating it
	// to ensure authenticity of the data
	currentBook, err := getBookById(book.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if (currentBook.Quantity == 0) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Book out of stock"})
		return
	}

	// Update the book instance with the new data
	update := bson.M{
		"$set": bson.M{
			"quantity": currentBook.Quantity - 1,
		},
	}

	// updating the book data
	filter := bson.M{"_id": currentBook.ID}
	_, updateErr := mongodb.BooksCollection.UpdateOne(context.Background(), filter, update)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the book data"})
		return
	}

	c.Status(http.StatusOK)
}