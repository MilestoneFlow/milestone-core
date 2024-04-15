package authorization

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetWorkspaceIDByUserIdentifier(workspaceCollection *mongo.Collection, userIdentifier string) (string, error) {
	projection := bson.D{{"_id", 1}}
	opts := options.FindOne().SetProjection(projection)

	var workspace bson.M
	err := workspaceCollection.FindOne(context.Background(), bson.M{"userIdentifiers": bson.M{"$in": []string{userIdentifier}}}, opts).Decode(&workspace)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	workspaceIdStr := workspace["_id"].(primitive.ObjectID).Hex()
	return workspaceIdStr, nil
}
