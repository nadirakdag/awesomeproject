package data

import (
	"awesomeProject/helpers"
	"awesomeProject/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// record repository interface
type RecordRepository interface {
	Get(filter *models.RecordFilter) ([]models.Record, error)
}

// record mongodb repository
type MongoRecordRepository struct {
	collection *mongo.Collection
}

// creates a new mongodb repository for record
func NewMongoRecordRepository(collection *mongo.Collection) *MongoRecordRepository {
	return &MongoRecordRepository{
		collection: collection,
	}
}

// implements Get method from RecordRepository for mongodb
// returns filtered records
func (mongoRepository *MongoRecordRepository) Get(filter *models.RecordFilter) ([]models.Record, error) {

	startDate, _ := time.Parse(helpers.DateFormat, filter.StartDate) // format checked on api
	endDate, _ := time.Parse(helpers.DateFormat, filter.EndDate)     // format checked on api

	pipe := []bson.M{
		{"$project": bson.M{"key": 1, "createdAt": 1, "totalCount": bson.M{"$sum": "$counts"}}},
		{"$match": bson.M{
			"totalCount": bson.M{"$gte": filter.MinCount, "$lte": filter.MaxCount},
			"createdAt":  bson.M{"$gte": startDate, "$lte": endDate},
		}},
	}

	cursor, err := mongoRepository.collection.Aggregate(context.TODO(), pipe)
	if err != nil {
		return nil, err
	}

	var result []models.Record
	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	return result, nil
}
