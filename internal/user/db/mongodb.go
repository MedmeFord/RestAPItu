package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/MedmeFord/RestAPItu/internal/user"
	"github.com/MedmeFord/RestAPItu/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user ")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}
	d.logger.Debug("convert insertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to conver objectId to Hex. oid: %s", oid)
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to ObjectIS. hex: %s", oid)
	}
	// mongo.getDatabase("test").getCOllection("docs").find({})
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			// TODO ErrEntityNotFound
			return u, fmt.Errorf("not found")
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(%s) from DB due to error: %v", id, err)
	}

	return u, nil
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users due to error: %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read documents from cursor, error: %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID) // -- Забираем айди и этот айди конвертируем из Hex в ojbjectID
	if err != nil {
		return fmt.Errorf("failed to conver to userID to ObjectId. Id=%s", user.ID)
	}

	filter := bson.M{"_id": objectID} // bson.M"{полеЭ: к чему равно} -- ищем определленого пользователя с определенныи айди

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v")
	}

	var updateUserObj bson.M // mongoDB принимает формат bson.M. Нам надо преоброзовать user байты в формат json
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("fsiled to unmarshal user bytes. eror: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute user query")
	}
	if result.MatchedCount == 0 { // MatchCount - проверяет наличие сущности. mongoDB сначала ищет сущность которую нужно обновить, а потом обновляет
		// TODO ErorEntity not found 404
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Mached %d documents and Modified %d document ", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id) // -- Забираем айди и этот айди конвертируем из Hex в ojbjectID
	if err != nil {
		return fmt.Errorf("failed to conver to userID to ObjectId. Id=%s", id)
	}

	filter := bson.M{"_id": objectID} // b

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Errorf("failed to execute to query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		// TODO ErorEntity not found 404
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, colletction string, logger *logging.Logger) user.Storage {

	return &db{
		collection: database.Collection(colletction),
		logger:     logger,
	}
}
