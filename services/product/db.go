package product

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductDatabase struct {
	productcollection *mongo.Collection
}

func NewProductDatabase(prodcoll *mongo.Collection) *ProductDatabase {
	return &ProductDatabase{
		productcollection: prodcoll,
	}
}
func (db *ProductDatabase) AddProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	_, err := db.productcollection.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	return nil
}
func (db *ProductDatabase) SearchProductByName(query string) (*[]Product, error) {
	var searchproduct []Product
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	searchresult, err := db.productcollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": query}})
	if err != nil {
		return nil, err
	}
	if err := searchresult.All(ctx, &searchproduct); err != nil {
		return nil, err
	}
	defer searchresult.Close(ctx)
	if err := searchresult.Err(); err != nil {
		return nil, err
	}
	return &searchproduct, nil

}
