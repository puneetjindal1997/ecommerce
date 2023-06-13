package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// Cart model `my_cart` table
type Cart struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Checkout  bool               `json:"checkout,omitempty" bson:"checkout"`
}
