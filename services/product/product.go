package product

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Product_ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Product_Name string             `json:"product_name"`
	Price        uint64             `json:"price"`
	Rating       uint8              `json:"rating"`
}
