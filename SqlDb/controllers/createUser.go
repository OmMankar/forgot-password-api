package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
)

type Create struct {
	Db *sql.DB
}

func (p *Create) CreateUser(w http.ResponseWriter, r *http.Request) {

	// Reading data from the request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Inserting the Data in Database
	query := "INSERT INTO user (emailId,password) VALUES (?, ?)"
	result, err := p.Db.ExecContext(context.Background(), query, user.EmailId, user.PassWord)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	InsertId, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	obj := map[string]interface{}{"InsertId": InsertId}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)

}
