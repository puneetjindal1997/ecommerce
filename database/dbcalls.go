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

func (mgr *manager) GetSingleProductById(id primitive.ObjectID, collectionName string) (product types.Product, err error) {
	filter := bson.D{{"_id", id}}
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)

	err = orgCollection.FindOne(context.TODO(), filter).Decode(&product)
	return product, err
}

func (mgr *manager) UpdateProduct(p types.Product, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"_id", p.ID}}
	update := bson.D{{"$set", p}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (mgr *manager) DeleteProduct(id primitive.ObjectID, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"_id", id}}

	_, err := orgCollection.DeleteOne(context.TODO(), filter)
	return err
}

func (mgr *manager) GetSingleAddress(id primitive.ObjectID, collectionName string) (address types.Address, err error) {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"user_id", id}}

	err = orgCollection.FindOne(context.TODO(), filter).Decode(&address)
	return address, err
}

func (mgr *manager) GetSingleUserByUserId(id primitive.ObjectID, collectionName string) (user types.User) {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"_id", id}}

	_ = orgCollection.FindOne(context.TODO(), filter).Decode(&user)
	return user
}

func (mgr *manager) GetCartObjectById(id primitive.ObjectID, collectionName string) (c types.Cart, err error) {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"_id", id}}

	err = orgCollection.FindOne(context.TODO(), filter).Decode(&c)
	return c, err
}

func (mgr *manager) UpdateUser(u types.User, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"_id", u.Id}}
	update := bson.D{{"$set", u}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
func (mgr *manager) UpdateCartToCheckOut(c types.Cart, collectionName string) error {
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{"_id", c.ID}}
	update := bson.D{{"$set", c}}

	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
