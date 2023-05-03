package utils

import (
	"config-generator/models"
	"context"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// TODO: move this to an env?
var mongoURI = "mongodb://localhost:27017"

const (
	databaseName = "AppData" // TODO: move this to env?
)

func CreateAppDataDB(appName string, version string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("Could not create new mongo client: %s", err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Could not connect to mongo client: %s", err)
	}

	// create database
	appData := models.DBAppData{
		Version: version,
		AppID:   GenerateUUID(),
	}

	DBRequest := models.CreateDBData{
		AppName: appName, // todo: change this to username
		Apps:    appData,
	}

	CreateDatabase(ctx, client, DBRequest)

	defer client.Disconnect(ctx)
}

func CreateDatabase(ctx context.Context, mongoClient *mongo.Client, DBRequest models.CreateDBData) {
	baseCollection := mongoClient.Database(databaseName).Collection(DBRequest.AppName)
	baseCollection.InsertOne(ctx, DBRequest)
}

func GenerateUUID() string {
	return "app-" + shortuuid.New()
	//return prefix + uuid.New().String()
}

func ReadFromDB(collectionName string) string {
	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Get a handle for the specified collection
	collection := client.Database(databaseName).Collection(collectionName)

	// Define a filter for the document with the specified AppName
	filter := bson.M{"appName": collectionName}

	// Define a variable to hold the resulting DB data
	var result models.CreateDBData

	// Retrieve the document from the collection
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Printf("Could not find from collection: %s", err)
	}

	// Return the AppID field for the retrieved document
	return result.Apps.AppID
}

// old version

//func ConnectMongoDB2() {
//	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ctx := context.Background()
//	err = client.Connect(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// create database
//	deploymentData := models.DeploymentData{
//		DeploymentName: deploymentName,
//		DeploymentID:   GenerateUUID(uuidPrefix),
//	}
//	appData := models.DBAppData{AppName: appName, Deployments: []models.DeploymentData{deploymentData}}
//
//	DBRequest := models.CreateDBData{
//		User: userName,
//		Apps: []models.DBAppData{appData},
//	}
//
//	CreateDatabase(ctx, client, DBRequest)
//
//	defer client.Disconnect(ctx)
//}

func ReadFromDB2(collectionName string, deploymentName string) string {
	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Get a handle to the collection
	collection := client.Database(databaseName).Collection(collectionName)

	projection := bson.D{
		{"deployments.deploymentid", 1},
		{"_id", 0},
	}

	cursor, err := collection.Find(context.TODO(), bson.D{{"deployments.deploymentname", deploymentName}}, options.Find().SetProjection(projection))

	//Iterate through the cursor and print the documents
	var deploymentId string
	for cursor.Next(ctx) {
		var result bson.M
		err = cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		deployments, isPrimitiveA := result["deployments"].(primitive.A)
		if isPrimitiveA {
			deployment, isBsonM := deployments[0].(bson.M)
			if isBsonM {
				deploymentId = deployment["deploymentid"].(string)
			}
		}
	}

	return deploymentId
}
