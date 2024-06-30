package mongoschemes

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID 		primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title	string `bson:"title" json:"title"`
	Author 	string `bson:"author" json:"author"`
	Quantity int   `bson:"quantity" json:"quantity"`
	PublishDate time.Time `bson:"publish_date" json:"publish_date"`
	Publisher string `bson:"publisher" json:"publisher"`
	PublisherId primitive.ObjectID `bson:"publisher_id" json:"publisher_id"`
	Rating int `bson:"rating" json:"rating"`
}