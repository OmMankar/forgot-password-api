# **User Authentication & Password Reset API (Go + MongoDB)**

## **Overview**

This is a RESTful API built using **Golang** and **MongoDB**, designed for user authentication and password reset functionality. The API allows users to create accounts, request OTPs for password reset, verify OTPs, and update their passwords securely.

## **Features**

- **User Registration**: Store user credentials in MongoDB.
- **Forgot Password Request**: Generate and send OTP via email.
- **OTP Verification**: Validate the OTP stored in cache.
- **Password Update**: Securely update the password after OTP verification.
- **Use Caching OTPs**: In-memory caching using `ttlcache` with a 3-minute expiration.
- **MongoDB Integration**: Stores user credentials and performs authentication.

## **Tech Stack**

- **Golang**
- **MongoDB**
- **Gorilla Mux** (Router)
- **ttlcache** (For OTP caching)
- **bcrypt** (For password hashing)
- **Gmail SMTP** (For sending OTPs via email)

## **Project Structure**

```
├── main.go            # Entry point of the application
├── controller/        # API handlers
│   ├── create.go      # Handles user creation
│   ├── forgot_pass.go # Handles OTP generation and password reset
├── database/          # Database connection setup
│   ├── connection.go  # MongoDB connection logic
├── models/            # Data models
│   ├── user.go        # User schema
├── template.html      # Email template for OTP
├── .env               # Environment variables (MongoDB URL, SMTP credentials)
└── README.md          # Documentation
```

## **Installation & Setup**

### **1. Clone the Repository**

```sh
git clone https://github.com/your-repo.git
cd your-repo
```

### **2. Install Dependencies**

```sh
go mod tidy
```

### **3. Set Up Environment Variables**

Create a `.env` file and add:

```ini
URL=mongodb+srv://your-mongo-url
SENDER_EMAIL=your-email@gmail.com
SENDER_PASSWORD=your-email-password
```

### **4. Run the Server**

```sh
go run main.go -port=3000
```

## **API Endpoints**

### **User Registration**

**POST** `/api/v1/`

```json
{
  "emailId": "user@example.com",
  "password": "securepassword"
}
```

### **Forgot Password Request (OTP Generation)**

**PUT** `/api/v1/forgot/password/request`

```json
{
  "emailId": "user@example.com"
}
```

### **OTP Verification**

**POST** `/api/v1/otp/check`

```json
{
  "emailId": "user@example.com",
  "Otp": 123456
}
```

### **Change Password**

**PUT** `/api/v1/new/password`

```json
{
  "emailId": "user@example.com",
  "password": "newsecurepassword"
}
```

## **Future Enhancements**

- Implement JWT authentication for enhanced security.
- Add rate limiting to prevent OTP abuse.
- Improve email security by using OAuth for SMTP authentication.

## **License**

This project is licensed under the **MIT License**.

