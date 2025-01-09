package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Create struct {
	Client *mongo.Client
}

func (p *Create) CreateUser(w http.ResponseWriter, r *http.Request) {
	collection := p.Client.Database("ForPass").Collection("User")

	// Reading data from the request body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}
	var user models.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}

	//Inserting the value in Database.
	doc := bson.M{"emailId": user.EmailId,
		"password": user.Password,
	}

	InsertOneResult, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		fmt.Println(err.Error())
	}

	obj := map[string]interface{}{"InsertId": InsertOneResult.InsertedID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)

}
