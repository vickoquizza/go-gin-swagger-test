package usecases

import (
	"context"
	"go-gin-swagger-test/app/db"
	"go-gin-swagger-test/core/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

const service_name = "gin-swagger-service"

type MongoConnection struct {
	*mongo.Collection
}

func NewMongoConnection() repository.Repository {
	collection, _ := db.NewMongoCollection("mongodb://root:root@localhost:27017/")

	return &MongoConnection{Collection: collection}
}

func (m *MongoConnection) InsertAccount(ctx context.Context, account db.AccountDTO) (db.Account, error) {
	PersistenceCurrentID++

	newCtx, span := otel.Tracer(service_name).Start(ctx, "Insert account document onto database")
	defer span.End()

	data := db.Account{
		Name: account.Name,
		ID:   PersistenceCurrentID,
	}

	_, err := m.InsertOne(newCtx, data)

	if err != nil {
		return db.Account{}, err
	}

	return data, nil
}

func (m *MongoConnection) GetAllAccounts(ctx context.Context) ([]db.Account, error) {
	var account db.Account
	var accounts []db.Account

	newCtx, span := otel.Tracer(service_name).Start(ctx, "Get all accounts documents")
	defer span.End()

	cursor, err := m.Find(ctx, bson.D{})

	if err != nil {
		defer cursor.Close(newCtx)
		return nil, err
	}

	for cursor.Next(newCtx) {
		err := cursor.Decode(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (m *MongoConnection) GetAccountById(ctx context.Context, id int) (db.Account, error) {
	var account db.Account

	newCtx, span := otel.Tracer(service_name).Start(ctx, "Get an account document by id")
	defer span.End()

	err := m.FindOne(newCtx, bson.D{{Key: "id", Value: id}}).Decode(&account)

	if err != nil {
		return db.Account{}, err
	}

	return account, nil
}

func (m *MongoConnection) UpdateAccountById(ctx context.Context, id int, accountUpdate db.AccountDTO) error {
	newCtx, span := otel.Tracer(service_name).Start(ctx, "Update an account document by id")
	defer span.End()

	filter := bson.D{{Key: "id", Value: id}}
	update_field := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: accountUpdate.Name}}}}

	_, err := m.UpdateOne(newCtx, filter, update_field)

	if err != nil {
		return err
	}

	return nil
}

func (m *MongoConnection) DeleteAccountById(ctx context.Context, id int) error {
	newCtx, span := otel.Tracer(service_name).Start(ctx, "Delete an account document by id")
	defer span.End()

	_, err := m.DeleteOne(newCtx, bson.D{{Key: "id", Value: id}})
	if err != nil {
		return err
	}

	return nil
}
