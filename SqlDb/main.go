package main

import (
	"flag"
	"fmt"
	"main/controllers"
	database "main/database"
	"net/http"
	"strconv"

	"github.com/ReneKroon/ttlcache"
	"github.com/gorilla/mux"
)

func main() {
	PORT := flag.Int("port", 3000, "Port of server on which our service will run")
	flag.Parse()

	//Connection to db.
	Db, err := database.Connect()
	if err != nil {
		fmt.Println("Error while connecting to database : ", err.Error())
		return
	}
	defer Db.Close()

	//Creating an instance of server.
	r := mux.NewRouter()

	//Creatin sub routes
	SubRouter := r.PathPrefix("/api/v1").Subrouter()

	//Creating cache of ttl/cache.
	Cache := ttlcache.NewCache()
	defer Cache.Close()

	f := controllers.ForPass{Db: Db, Cache: Cache}
	c := controllers.Create{Db: Db}
	SubRouter.HandleFunc("/", c.CreateUser).Methods("POST", "OPTIONS")
	SubRouter.HandleFunc("/forgot/password/request", f.RequestForPass).Methods("PUT", "OPTIONS")
	SubRouter.HandleFunc("/otp/check", f.CheckOtp).Methods("POST", "OPTIONS")
	SubRouter.HandleFunc("/new/password", f.ChangePass).Methods("PUT", "OPTIONS")

	//Default Route.
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Server"))

	})

	//Server listen on the port number.
	address := ":" + strconv.Itoa(*PORT)
	fmt.Printf("Server to started On PORT Number : %v ", *PORT)
	http.ListenAndServe(address, r)

}
