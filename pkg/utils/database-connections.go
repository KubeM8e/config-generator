package utils

import (
	"config-generator/models"
	"context"
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
	databaseName   = "AppStructure" // TODO: move this to env?
	deploymentName = "Deployment-1" // TODO: get this from user or auto generate?
	userName       = "Jon Doe"      // TODO: move this to env?
	uuidPrefix     = "app-"         // TODO: move this to env?
	appName        = "my-demo-app"
)

func ConnectMongoDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// create database
	deploymentData := models.DeploymentData{
		DeploymentName: deploymentName,
		DeploymentID:   GenerateUUID(uuidPrefix),
	}
	appData := models.DBAppData{AppName: appName, Deployments: []models.DeploymentData{deploymentData}}

	DBRequest := models.CreateDBData{
		User: userName,
		Apps: []models.DBAppData{appData},
	}

	CreateDatabase(ctx, client, DBRequest)

	defer client.Disconnect(ctx)
}

func CreateDatabase(ctx context.Context, mongoClient *mongo.Client, DBRequest models.CreateDBData) {
	baseCollection := mongoClient.Database(databaseName).Collection(DBRequest.User)
	for _, app := range DBRequest.Apps {
		baseCollection.InsertOne(ctx, app)
	}

}

func GenerateUUID(prefix string) string {
	return prefix + shortuuid.New()
	//return prefix + uuid.New().String()
}

func ReadFromDB(collectionName string, deploymentName string) string {
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
