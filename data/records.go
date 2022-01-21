package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Record struct {
	Key        string    `json:"key" bson:"key"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	TotalCount int       `json:"totalCount" bson:"totalCount"`
}

type RecordFilter struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type RecordRepository interface {
	Get(filter *RecordFilter) []Record
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

func (mongoRepository *MongoRecordRepository) Get(filter *RecordFilter) []Record {

	startDate, err := time.Parse("2006-01-02", filter.StartDate)
	if err != nil {
		log.Panicln("start date parse failed", filter.StartDate)
		return nil
	}

	endDate, err := time.Parse("2006-01-02", filter.EndDate)
	if err != nil {
		log.Panicln("end date parse failed")
		return nil
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
		panic(err)
	}

	var result []Record
	if err := cursor.All(mongoRepository.context, &result); err != nil {
		log.Fatalln(err)
	}

	return result
}
