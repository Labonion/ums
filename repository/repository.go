package repository

import (
	"context"
	"log"
	"markie-backend/database"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AbstractRepository struct {
	collection *mongo.Collection
	name       string
	client     *mongo.Client
}

func Repository(collectionName string) *AbstractRepository {
	collection := database.DB.Database(os.Getenv("DATABASE")).Collection(collectionName)
	return &AbstractRepository{
		collection: collection,
		name:       collection.Name(),
		client:     database.DB,
	}
}

func (repository *AbstractRepository) Aggregate(filter *bson.M, additionalStages ...interface{}) (*mongo.Cursor, error) {
	session, sessionErr := repository.client.StartSession()
	if sessionErr != nil {
		log.Fatal(sessionErr)
		return nil, sessionErr
	}
	defer session.EndSession(context.Background())
	transactionErr := session.StartTransaction()
	if transactionErr != nil {
		log.Fatal(transactionErr)
		return nil, transactionErr
	}

	pipeline := bson.A{
		bson.M{"$match": *filter},
	}
	pipeline = append(pipeline, additionalStages...)

	cur, err := repository.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
		session.AbortTransaction(context.Background())
		return nil, err
	}
	return cur, nil
}

func (repository *AbstractRepository) FindAll(filter *bson.M) (*mongo.Cursor, error) {
	cur, err := repository.collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cur, nil
}

func (repository *AbstractRepository) InsertOne(doc interface{}) (*mongo.InsertOneResult, error) {
	session, sessionErr := repository.client.StartSession()
	if sessionErr != nil {
		log.Fatal(sessionErr)
		return nil, sessionErr
	}
	defer session.EndSession(context.Background())
	transactionErr := session.StartTransaction()
	if transactionErr != nil {
		log.Fatal(transactionErr)
		return nil, transactionErr
	}
	cur, err := repository.collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Fatal(err)
		session.AbortTransaction(context.Background())
		return nil, err
	}
	return cur, nil
}

func (repository *AbstractRepository) UpdateById(Id primitive.ObjectID, updateDoc interface{}) (*mongo.UpdateResult, error) {
	session, sessionErr := repository.client.StartSession()
	if sessionErr != nil {
		log.Fatal(sessionErr)
		return nil, sessionErr
	}
	defer session.EndSession(context.Background())

	transactionErr := session.StartTransaction()
	if transactionErr != nil {
		log.Fatal(transactionErr)
		return nil, transactionErr
	}

	cur, err := repository.collection.UpdateOne(context.Background(), bson.M{"_id": Id}, updateDoc)
	if err != nil {
		log.Fatal(err)
		session.AbortTransaction(context.Background())
		return nil, err
	}

	if commitErr := session.CommitTransaction(context.Background()); commitErr != nil {
		log.Fatal(commitErr)
		return nil, commitErr
	}

	return cur, nil
}

func (repository *AbstractRepository) FindOne(filter *bson.M) *mongo.SingleResult {
	cur := repository.collection.FindOne(context.TODO(), filter)
	if cur == nil {
		return nil
	}
	return cur
}
