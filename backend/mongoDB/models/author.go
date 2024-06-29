package mongoschemes

import "go.mongodb.org/mongo-driver/bson/primitive"

type Author struct {
	ID 		primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username	string `bson:"username" json:"username"`
	Password 	string `bson:"password" json:"password"`
	PublishedBooks []primitive.ObjectID `bson:"published_books"`
}