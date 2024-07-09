package mongoapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mongodb "github.com/roeeattias/Book-Store/mongoDB/database"
	mongoschemes "github.com/roeeattias/Book-Store/mongoDB/models"
	requestsDataStructures "github.com/roeeattias/Book-Store/mongoDB/requestsStructs"
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

	base64ImageData, err := encodeImageToBase64(author.ImageUrl)
	if (err != nil) {
		c.Status(http.StatusInternalServerError)
		return
	}
	
	// generating a jwt token
	token, err := generateJWTtoken(credentials.Username, author.ID.Hex())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": author.ID,
		"username": author.Username,
		"publishedBooks": author.PublishedBooks,
		"image_url": base64ImageData,
	}

	// setting an authorization cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600 * 24 * 30, "", "", false, true)
	c.JSON(http.StatusOK, response)
}

func SignUp(c *gin.Context) {
	var newAuthor mongoschemes.Author

	// pulling author data from the request to the object
	if err := c.BindJSON(&newAuthor); err != nil {
		return
	}
	
	// checking if the author exists in the database
	author := getAuthorByName(newAuthor.Username)
	if author.Username != "" {
		fmt.Println("here1")
		c.Status(http.StatusInternalServerError)
		return
	}

	decodedImageData, base64ImageData := decodeImageFromBase64(newAuthor.ImageUrl)
	if (decodedImageData == nil) {
		fmt.Println("here2")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Determine the file type using http.DetectContentType
	fileType := http.DetectContentType(decodedImageData)
	filePath := "profileImages/" + newAuthor.Username + "." + strings.Split(fileType, "/")[1]
	
    // Write to file
    if err := os.WriteFile(filePath, decodedImageData, 0644); err != nil {
		fmt.Println("here3")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Faild to upload profile image"})
		return
    }

	// hashing the password (Bcrypt, 12 factor)
	password := newAuthor.Password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
    newAuthor.Password = string(bytes)
	newAuthor.PublishedBooks = []primitive.ObjectID{}
	newAuthor.ImageUrl = filePath

	// Creating a new user document
	authorInstance, err := mongodb.AuthorCollection.InsertOne(context.Background(), newAuthor)
	if err != nil {
		fmt.Println("here4")
		c.Status(http.StatusInternalServerError)
		return
	}

	// generating JWT token
	token, err := generateJWTtoken(newAuthor.Username, authorInstance.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		fmt.Println("here5")
		c.Status(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": authorInstance.InsertedID,
		"username": newAuthor.Username,
		"publishedBooks": newAuthor.PublishedBooks,
		"image_url": "data:image/jpeg;base64," + base64ImageData,
	}
	
	// setting the authorization cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600 * 24 * 30, "", "", false, true)
	c.JSON(http.StatusCreated, response)
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
		c.Status(http.StatusNetworkAuthenticationRequired)
        return
	}
	
	decodedImageData, base64ImageData := decodeImageFromBase64(newBook.ImageUrl)
	if (decodedImageData == nil) {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Determine the file type using http.DetectContentType
	fileType := http.DetectContentType(decodedImageData)
	filePath := "bookImages/" + generateImageIdentifier() + "." + strings.Split(fileType, "/")[1]
	
    // Write to file
    if err := os.WriteFile(filePath, decodedImageData, 0644); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Faild to upload profile image"})
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
	newBook.ImageUrl = filePath

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author's books"})
		return
	}
	
	newBookId, ok := bookInstance.InsertedID.(primitive.ObjectID)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	newBook.ID = newBookId
	newBook.ImageUrl = "data:image/jpeg;base64," + base64ImageData;
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
		base64ImageData, err := encodeImageToBase64(book.ImageUrl)
		if (err != nil) {

			c.Status(http.StatusInternalServerError)
			return
		}
		book.ImageUrl = base64ImageData
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

	// Update the authors collection to remove the book from the author's books array
	update := bson.M{
		"$pull": bson.M{
			"published_books": currentBook.ID, // Assuming bookId is the ID of the book you want to delete
		},
	}

	filter := bson.M{"_id": userObjectId}

	_, updateErr := mongodb.AuthorCollection.UpdateOne(context.Background(), filter, update)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete books from author"})
		return
	}


	// deleting the book image from the backend local storage
	fmt.Println(currentBook.ImageUrl.(string))
	os.Remove(currentBook.ImageUrl.(string))
    

	// Define the filter to match the ObjectId
    filter = bson.M{"_id": bookToDelete.ID}

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

func GetAuthors(c *gin.Context) {
	query := c.PostForm("query")
		
	// Define the filter to find authors whose usernames start with the query
	filter := bson.M{"username": bson.M{"$regex": "^" + query, "$options": "i"}}

	// Use Find method with filter to get a cursor for the matching documents
	cursor, err := mongodb.AuthorCollection.Find(context.Background(), filter)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Prepare a slice to hold the results
	var authors []mongoschemes.Author

	// Iterate through the cursor and decode directly into the response slice
	for cursor.Next(context.Background()) {
		var author mongoschemes.Author
		if err := cursor.Decode(&author); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		base64ImageData, err := encodeImageToBase64(author.ImageUrl)
		if (err != nil) {
			c.Status(http.StatusInternalServerError)
			return
		}
		author.ImageUrl = base64ImageData
		authors = append(authors, author)
	}

	// Check for errors after iterating through the cursor
	if err := cursor.Err(); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(authors) > 0 {
		// Send the response
		c.JSON(http.StatusOK, authors)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func GetAuthorBooks(c *gin.Context) {
	var books requestsDataStructures.AuthorsPublishedBooks
	
	// Bind JSON payload to the struct
	if err := c.BindJSON(&books); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Convert the array of string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, len(books.Books))
	for i, idStr := range books.Books {
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID format"})
			return
		}
		objectIDs[i] = objectID
	}
	
	// Use the $in operator to find matching documents
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := mongodb.BooksCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding books"})
		return
	}
	defer cursor.Close(context.Background())

	var booksFound []mongoschemes.Book
	for cursor.Next(context.Background()) {
		var book mongoschemes.Book
		if err := cursor.Decode(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding book"})
			return
		}
		base64ImageData, err := encodeImageToBase64(book.ImageUrl)
		if (err != nil) {

			c.Status(http.StatusInternalServerError)
			return
		}
		book.ImageUrl = base64ImageData
		booksFound = append(booksFound, book)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	// Return the matching books as a JSON response
	c.JSON(http.StatusOK, booksFound)
}
