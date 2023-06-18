package database

import (
	"context"
	"ecommerce/constant"
	"ecommerce/types"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mgr *manager) Insert(data interface{}, collectionName string) (interface{}, error) {
	log.Println(data)
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	// insert the bson object using InsertOne()
	result, err := orgCollection.InsertOne(context.TODO(), data)
	// check for errors in the insertion
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (mgr *manager) GetSingleRecordByEmail(email string, collectionName string) *types.Verification {
	resp := &types.Verification{}
	filter := bson.D{{"email", email}}
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	_ = orgCollection.FindOne(context.TODO(), filter).Decode(&resp)
	fmt.Println(resp)
	return resp

}

func (mgr *manager) UpdateVerification(data types.Verification, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	filter := bson.D{{"email", data.Email}}
	update := bson.D{{"$set", data}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)

	return err

}

func (mgr *manager) UpdateEmailVerifiedStatus(req types.Verification, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"email", req.Email}}
	update := bson.D{{"$set", req}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)

	return err
}

// Get single user from db
func (mgr *manager) GetSingleRecordByEmailForUser(email, collectionName string) types.User {
	resp := types.User{}
	filter := bson.D{{"email", email}}
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	_ = orgCollection.FindOne(context.TODO(), filter).Decode(&resp)
	fmt.Println(resp)
	return resp
}

func (mgr *manager) GetListProducts(page, limit, offset int, collectionName string) (products []types.Product, count int64, err error) {
	skip := ((page - 1) * limit)
	if offset > 0 {
		skip = offset
	}
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cur, err := orgCollection.Find(context.TODO(), bson.M{}, findOptions)
	err = cur.All(context.TODO(), &products)
	itemCount, err := orgCollection.CountDocuments(context.TODO(), bson.M{})
	return products, itemCount, err
}

func (mgr *manager) SearchProduct(page, limit, offset int, search string, collectionName string) (products []types.Product, count int64, err error) {
	skip := ((page - 1) * limit)
	if offset > 0 {
		skip = offset
	}

	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	searchFilter := bson.M{}
	if len(search) >= 3 {
		searchFilter["$or"] = []bson.M{
			{"name": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
			{"description": primitive.Regex{Pattern: ".*" + search + ".*", Options: "i"}},
		}
	}
	cur, err := orgCollection.Find(context.TODO(), searchFilter, findOptions)
	cur.All(context.TODO(), &products)
	count, err = orgCollection.CountDocuments(context.TODO(), searchFilter)
	return products, count, err
}
