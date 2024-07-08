package mongoschemes

import "go.mongodb.org/mongo-driver/bson/primitive"

type Author struct {
	ID 		primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username	string `bson:"username" json:"username"`
	Password 	string `bson:"password" json:"password"`
	ImageUrl interface{} `bson:"image_url" json:"image_url"`
	PublishedBooks []primitive.ObjectID `bson:"published_books" json:"publishedBooks"`
}