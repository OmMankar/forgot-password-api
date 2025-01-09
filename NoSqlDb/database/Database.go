package database

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connection() (Client *mongo.Client, err error) {

	//Load .env file variable in environment variables.
	godotenv.Load(".env")
	url := os.Getenv("URL")

	//Connecting to Mongodb.
	Client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(url))
	if err != nil {
		fmt.Println("Error in connecting with Db")

		return nil, err
	}
	fmt.Println("Database connection successful")

	return Client, nil
}
