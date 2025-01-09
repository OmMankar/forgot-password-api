package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"math/rand/v2"
	"net/http"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type ForPass struct {
	Db    *sql.DB
	Cache *ttlcache.Cache
}

func SendEmail(to []string, otp int) error {
	//sender data.
	godotenv.Load(".env")
	from := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SENDER_PASSWORD")

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message
	t, err := template.ParseFiles("./template.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	var body bytes.Buffer
	// this we need to write when using an html file.
	mimeHeaders := "MIME-version: 1.0;\nContent-Type:text/html;charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Forgot Password Request \n%s \n\n", mimeHeaders)))
	t.Execute(&body, struct {
		OTP int
	}{OTP: otp})

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {

		return err
	}
	fmt.Println("Email Sent Successfully!")
	return nil
}

func (p *ForPass) RequestForPass(w http.ResponseWriter, r *http.Request) {

	//Parsing Json data from Request Body.
	obj := models.User{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		fmt.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validating the email ID.
	if obj.EmailId == "" {
		response := map[string]string{"message": "Invalid Email Id Enterd", "success": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Checking whether the email ID exists in the database.
	query := "SELECT * FROM user WHERE emailId=?"
	rows, err := p.Db.QueryContext(context.Background(), query, obj.EmailId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User Does Not Exist")
			response := map[string]string{"message": "No Such user Present in Db", "success": "false"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		} else {
			fmt.Println("Error While checking", err.Error())
			response := map[string]string{"message": "Error occured while finding the user in Db", "success": "false"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

	}
	rows.Close()

	// Generating a 6-digit OTP.
	otp := rand.IntN(800000) + 100000
	fmt.Println(otp)

	// Caching the OTP in memory with a 3-minute expiration time.
	p.Cache.SetWithTTL(obj.EmailId, otp, (3)*time.Minute)

	// Sending the OTP to the client's email.
	err = SendEmail([]string{obj.EmailId}, otp)
	if err != nil {
		fmt.Println(err.Error())
		response := map[string]string{"message": "Unable to generate OTP", "success": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Sending the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{ Message string }{Message: "Successfuly send otp on email"})

}

func (p *ForPass) CheckOtp(w http.ResponseWriter, r *http.Request) {

	// Pasing json data from the request body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		response := map[string]string{"message": "Error occured in reading Request Body", "success": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	obj := struct {
		Otp     int
		EmailId string
	}{}
	json.Unmarshal(data, &obj)

	// Received the email ID as a string.
	clientOtp := obj.Otp
	fmt.Println("clientOtp : ", clientOtp)

	// Validating the email ID.
	if clientOtp <= 99999 || clientOtp >= 1000000 {
		response := map[string]string{"message": "Invalid OTP Entered", "success": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	generatedOtp, status := p.Cache.Get(obj.EmailId)
	fmt.Printf("genrate :%v , status:%v", generatedOtp, status)
	if status {
		if generatedOtp == clientOtp {
			response := map[string]string{"message": "OTP Verified", "success": "true"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		} else {
			response := map[string]string{"message": "Incorrect OTP", "success": "false"}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	//OTP got exired.
	response := map[string]string{"message": "Generate OTP once again", "success": "false"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)

}
func (p *ForPass) ChangePass(w http.ResponseWriter, r *http.Request) {

	// Parsing JSON data from the request body.
	var obj models.User
	json.NewDecoder(r.Body).Decode(&obj)

	// Received the Password as a string.
	newPass := obj.PassWord

	// Validating the password.
	if newPass == "" {
		response := map[string]string{"message": "Password is empty entered", "success": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}

	// Adding the encrypted password to the database.
	hash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		response := map[string]string{"message": "Error occured while generating hash code of User Input Password", "success": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Updating the password in the database.
	query := "UPDATE user SET password=? WHERE emailId = ?"
	_, err = p.Db.ExecContext(context.Background(), query, hash, obj.EmailId)
	if err != nil {
		fmt.Println(err.Error())
		response := map[string]string{"message": "Error occured while updating Db", "success": "false"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Sending Response.
	response := map[string]string{"message": "Password updated successfully", "success": "true"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
