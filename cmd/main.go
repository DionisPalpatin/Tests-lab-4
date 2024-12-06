package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
)

var verificationCodes = make(map[string]string)
var userPassword string
var userEmail string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	userPassword = os.Getenv("USER_PASSWORD")
	userEmail = os.Getenv("USER_EMAIL")

	router := gin.Default()
	router.POST("/status", statusHandler)
	router.POST("/login", loginHandler)
	router.POST("/login/verify", verifyLoginHandler)
	router.POST("/password/reset", resetPasswordHandler)
	router.POST("/password/reset/verify", verifyResetHandler)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}

func statusHandler(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func loginHandler(ctx *gin.Context) {
	var request struct {
		Email    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if request.Password != userPassword || request.Email != userEmail {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
	}

	code := rand.IntN(123456789)
	codeStr := strconv.Itoa(code)
	verificationCodes[request.Email] = codeStr
	sendEmail(codeStr)

	ctx.JSON(http.StatusOK, gin.H{"massage": "Auth email send successfully"})
}

func verifyLoginHandler(ctx *gin.Context) {
	var request struct {
		Email string `json:"login" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	code, ok := verificationCodes[request.Email]
	if !ok || request.Code != code {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong code"})
	}

	delete(verificationCodes, request.Email)

	ctx.JSON(http.StatusOK, gin.H{"massage": "Auth is made successfully"})
}

func resetPasswordHandler(ctx *gin.Context) {
	var request struct {
		Email    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if request.Password != userPassword || request.Email != userEmail {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
	}

	code := rand.IntN(123456789)
	codeStr := strconv.Itoa(code)
	if os.Getenv("THIS_IS_TEST") != "" {
		codeStr = os.Getenv("TEST_PASSWORD")
	}
	verificationCodes[request.Email] = codeStr
	sendEmail(codeStr)

	ctx.JSON(http.StatusOK, gin.H{"massage": "Reset password code send successfully"})
}

func verifyResetHandler(ctx *gin.Context) {
	var request struct {
		Email       string `json:"login" binding:"required"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	code, ok := verificationCodes[request.Email]
	if !ok || request.Code != code {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong code"})
	}

	delete(verificationCodes, request.Email)
	userPassword = request.NewPassword

	ctx.JSON(http.StatusOK, gin.H{"massage": "Password is changed successfully"})
}

func sendEmail(code string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	senderEmail := os.Getenv("SENDER_EMAIL_ADDRESS")
	senderPassword := os.Getenv("SENDER_EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_SERVER")
	smtpPort := 587
	userEmail := os.Getenv("USER_EMAIL")

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Error sending verification code: %v", err)
	}
}
