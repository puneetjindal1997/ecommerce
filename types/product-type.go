package types

import "go.mongodb.org/mongo-driver/bson/primitive"

// Product model `products` table
type Product struct {
	ID          primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Price       float64                `json:"price" bson:"price"`
	ImageUrl    string                 `json:"image_url" bson:"image_url"`
	MetaInfo    map[string]interface{} `json:"meta_info,omitempty" bson:"meta_info,omitempty"`
	CreatedAt   int64                  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64                  `json:"updated_at" bson:"updated_at"`
}

type ProductClient struct {
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Price       float64                `json:"price" bson:"price"`
	ImageUrl    string                 `json:"image_url" bson:"image_url"`
	MetaInfo    map[string]interface{} `json:"meta_info,omitempty" bson:"meta_info,omitempty"`
}

type UpdateProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
}
