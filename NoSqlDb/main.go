package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"main/controller"
	"main/database"
	"net/http"

	"github.com/ReneKroon/ttlcache"
	"github.com/gorilla/mux"
)

func main() {
	//Taking PORT as input.
	PORT := flag.Int("port", 3000, "Port Number ")
	flag.Parse()

	//Connecting database.
	Client, err := database.Connection()
	if err != nil {
		fmt.Println(err.Error())
	}

	//Creating cache of ttl/cache.
	cache := ttlcache.NewCache()
	defer cache.Close()

	//Creating an instance of server.
	r := mux.NewRouter()

	//Mounting the api path with version.
	subRouter := r.PathPrefix("/api/v1").Subrouter()

	f := controller.ForPass{Client: Client, Cache: cache}
	c := controller.Create{Client: Client}
	subRouter.HandleFunc("/", c.CreateUser).Methods("POST", "OPTIONS")
	subRouter.HandleFunc("/forgot/password/request", f.RequestForPass).Methods("PUT", "OPTIONS")
	subRouter.HandleFunc("/otp/check", f.CheckOtp).Methods("POST", "OPTIONS")
	subRouter.HandleFunc("/new/password", f.ChangePass).Methods("PUT", "OPTIONS")

	//Default Route.
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "This is the home page"})
	}).Methods("GET", "OPTION")

	Url := fmt.Sprintf(":%d", *PORT)
	fmt.Printf("Server Started at port %d\n", *PORT)
	http.ListenAndServe(Url, r)
}
