package data

import (
	"awesomeProject/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecordRepository interface {
	Get(filter *models.RecordFilter) ([]models.Record, error)
}

type MongoRecordRepository struct {
	collection *mongo.Collection
	context    context.Context
}

func NewMongoRecordRepository(collection *mongo.Collection, context context.Context) *MongoRecordRepository {
	return &MongoRecordRepository{
		collection: collection,
		context:    context,
	}
}

var ErrStartDateFormatInvalid = errors.New("start date format is invalid, is should be YYYY-MM-DD")
var ErrEndDateFormatInvalid = errors.New("end date format is invalid, is should be YYYY-MM-DD")

func (mongoRepository *MongoRecordRepository) Get(filter *models.RecordFilter) ([]models.Record, error) {

	startDate, err := time.Parse("2006-01-02", filter.StartDate)
	if err != nil {
		return nil, ErrStartDateFormatInvalid
	}

	endDate, err := time.Parse("2006-01-02", filter.EndDate)
	if err != nil {
		return nil, ErrEndDateFormatInvalid
	}

	pipe := []bson.M{
		{"$project": bson.M{"key": 1, "createdAt": 1, "totalCount": bson.M{"$sum": "$counts"}}},
		{"$match": bson.M{
			"totalCount": bson.M{"$gte": filter.MinCount, "$lte": filter.MaxCount},
			"createdAt":  bson.M{"$gte": startDate, "$lte": endDate},
		}},
	}

	cursor, err := mongoRepository.collection.Aggregate(mongoRepository.context, pipe)
	if err != nil {
		return nil, err
	}

	var result []models.Record
	if err := cursor.All(mongoRepository.context, &result); err != nil {
		return nil, err
	}

	return result, nil
}
