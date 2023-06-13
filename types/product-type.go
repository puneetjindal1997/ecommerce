package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product model `products` table
type Product struct {
	ID          primitive.ObjectID     `json:"_id" bson:"_id"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Price       float64                `json:"price" bson:"price"`
	ImageUrl    string                 `json:"image_url" bson:"image_url"`
	MetaInfo    map[string]interface{} `json:"meta_info,omitempty" bson:"meta_info,omitempty"`
}
