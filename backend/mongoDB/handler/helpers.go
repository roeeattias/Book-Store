package mongoapi

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	mongodb "github.com/roeeattias/Book-Store/mongoDB/database"
	mongoschemes "github.com/roeeattias/Book-Store/mongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func getAuthorByName(username string) mongoschemes.Author {
	var user mongoschemes.Author
	filter := bson.D{{Key: "username", Value: username}}
	mongodb.AuthorCollection.FindOne(context.Background(), filter).Decode(&user)
	return user
}

func generateJWTtoken(username string, id string) (string, error) {
	// Loading the jwt secret key used for generateing the token
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtKey) == 0 {
		return "", fmt.Errorf("JWT_SECRET_KEY not set or empty")
	}

	// creating the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id": id,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "BookStore",
		"aud":      "Authors",
	})

	// Signing the jwt token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	// Loading the secret key from env file
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("JWT_SECRET_KEY not set or empty")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing error: %v", err)
	}

	// Check if the token is valid and validate standard claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validate the issuer
		if claims["iss"] != "BookStore" {
			return nil, fmt.Errorf("invalid issuer")
		}
		// Validate the audience
		if claims["aud"] != "Authors" {
			return nil, fmt.Errorf("invalid audience")
		}
		// Validate the expiration time
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, fmt.Errorf("token expired")
		}
		return token, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func Middleware(c *gin.Context) {
	// Retrieve the token from the cookie
	jwtToken, err := c.Cookie("Authorization")

	if err != nil {
		c.Status(http.StatusSeeOther)
		c.Abort()
		return
	}

	// Verify the token
	token, err := verifyToken(jwtToken)
	if err != nil {
		c.Status(http.StatusSeeOther)
		c.Abort()
		return
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.Status(http.StatusSeeOther)
		c.Abort()
		return
	}

	// Extract the username from the claims
	username, ok := claims["username"].(string)
	if !ok {
		c.Status(http.StatusSeeOther)
		c.Abort()
		return
	}
	
	// Extract the user id from the claims
	userId, ok := claims["user_id"].(string)
	if !ok {
		c.Status(http.StatusSeeOther)
		c.Abort()
		return
	}

	c.Set("username", username)
	c.Set("userId", userId)

	// Continue with the next middleware or route handler
	c.Next()
}

func getBookById(id primitive.ObjectID) (mongoschemes.Book, error) {
	// Fetch the document by ObjectID
    var book mongoschemes.Book
    err := mongodb.BooksCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&book)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return mongoschemes.Book{}, errors.New("BOOK NOT FOUND")
        } else {
            return mongoschemes.Book{}, err
        }
    } else {
        return book, nil
    }
}

func encodeImageToBase64(imagePath interface{}) (string, error) {
	imageUrl, ok := imagePath.(string)
	if (!ok) {
		return "", errors.New("could not load image")
	}

	// Read image file
	imageFile, err := os.ReadFile(imageUrl)

	if err != nil {
		return "", err
	}

	// Encode image data to base64
	base64ImageData := base64.StdEncoding.EncodeToString(imageFile)
	base64ImageData = "data:image/jpeg;base64," + base64ImageData
	
	return base64ImageData, nil
}

func decodeImageFromBase64(imagePath interface{}) ([]byte, string) {
	// Type assertion to convert interface{} to map[string]interface{}
	fmt.Println(imagePath)
	dataMap, ok := imagePath.(map[string]interface{});
	if !ok {
		fmt.Println("here1")
		return nil, ""
	}

	// Access the 'dataUrl' key if it exists
	dataUrl, ok := dataMap["dataUrl"].(string);
	if !ok {
		fmt.Println("here2")
		return nil, ""
	}

	// Extract the base64 encoded data
    base64ImageData := strings.Split(dataUrl, ";base64,")[1]

    // Decode base64 string to byte slice
    decoded, err := base64.StdEncoding.DecodeString(base64ImageData)
    if err != nil {
		fmt.Println("here3")
		return nil, ""
    }

	return decoded, base64ImageData
}