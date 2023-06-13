package database

import (
	"context"
	"ecommerce/constant"
	"ecommerce/types"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func (mgr *manager) Insert(data interface{}, collectionName string) error {
	log.Println(data)
	orgCollection := mgr.connection.Database(constant.Database).Collection(collectionName)
	// insert the bson object using InsertOne()
	result, err := orgCollection.InsertOne(context.TODO(), data)
	// check for errors in the insertion
	if err != nil {
		return err
	}
	// display the id of the newly inserted object
	log.Println(result.InsertedID)

	return nil
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
