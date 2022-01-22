package data

import (
	"awesomeProject/db"
	"awesomeProject/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// record repository interface
type RecordRepository interface {
	Get(filter *models.RecordFilter) ([]models.Record, error)
}

// record mongodb repository
type MongoRecordRepository struct {
	connection *db.MongoDbConnection
}

// creates a new mongodb repository for record
func NewMongoRecordRepository(conection *db.MongoDbConnection) *MongoRecordRepository {
	return &MongoRecordRepository{
		connection: conection,
	}
}

var ErrStartDateFormatInvalid = errors.New("start date format is invalid, is should be YYYY-MM-DD")
var ErrEndDateFormatInvalid = errors.New("end date format is invalid, is should be YYYY-MM-DD")

const (
	dateFormat string = "2006-01-02"
	database   string = "getir-case-study"
	collection string = "records"
)

// implements Get method from RecordRepository for mongodb
// returns filtered records
// if filter object start date format is invalid then returns ErrStartDateFormatInvalid
// if filter object end date format is invalid then returns ErrEndDateFormatInvalid
func (mongoRepository *MongoRecordRepository) Get(filter *models.RecordFilter) ([]models.Record, error) {

	startDate, err := time.Parse(dateFormat, filter.StartDate)
	if err != nil {
		return nil, ErrStartDateFormatInvalid
	}

	endDate, err := time.Parse(dateFormat, filter.EndDate)
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

	if err := mongoRepository.connection.Client.Connect(context.TODO()); err != nil {
		return nil, err
	}

	if err := mongoRepository.connection.Client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	collection := mongoRepository.connection.Client.Database(database).Collection(collection)
	cursor, err := collection.Aggregate(context.TODO(), pipe)
	if err != nil {
		return nil, err
	}

	var result []models.Record
	if err := cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	if err := mongoRepository.connection.Client.Disconnect(context.TODO()); err != nil {
		return nil, err
	}

	return result, nil
}
