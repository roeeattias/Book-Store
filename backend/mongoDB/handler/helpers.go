package mongoapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	mongodb "github.com/roeeattias/Book-Store/mongoDB/database"
	mongoschemes "github.com/roeeattias/Book-Store/mongoDB/models"
	"go.mongodb.org/mongo-driver/bson"
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
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	// Verify the token
	token, err := verifyToken(jwtToken)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	// Extract the username from the claims
	username, ok := claims["username"].(string)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}
	
	// Extract the user id from the claims
	userId, ok := claims["user_id"].(string)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}
	
	c.Set("username", username)
	c.Set("userId", userId)

	// Continue with the next middleware or route handler
	c.Next()
}