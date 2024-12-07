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
	router.GET("/status", statusHandler)
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Password != userPassword || request.Email != userEmail {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	code := rand.IntN(123456789)
	codeStr := strconv.Itoa(code)
	if os.Getenv("THIS_IS_TEST") != "" {
		codeStr = os.Getenv("TEST_CODE")
	}

	verificationCodes[request.Email] = codeStr

	if err := sendEmail(codeStr); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Auth email send successfully"})
}

func verifyLoginHandler(ctx *gin.Context) {
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, ok := verificationCodes[request.Email]
	if !ok || request.Code != code {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong code"})
		return
	}

	delete(verificationCodes, request.Email)

	ctx.JSON(http.StatusOK, gin.H{"message": "Auth is made successfully"})
}

func resetPasswordHandler(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Password != userPassword || request.Email != userEmail {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	code := rand.IntN(123456789)
	codeStr := strconv.Itoa(code)
	if os.Getenv("THIS_IS_TEST") != "" {
		codeStr = os.Getenv("TEST_CODE")
	}

	verificationCodes[request.Email] = codeStr

	if err := sendEmail(codeStr); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reset password code send successfully"})
}

func verifyResetHandler(ctx *gin.Context) {
	var request struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, ok := verificationCodes[request.Email]
	if !ok || request.Code != code {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong code"})
		return
	}

	delete(verificationCodes, request.Email)
	userPassword = request.NewPassword

	ctx.JSON(http.StatusOK, gin.H{"message": "Password is changed successfully"})
}

func sendEmail(code string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	senderName := os.Getenv("SENDER_LOGIN")
	senderPassword := os.Getenv("SENDER_EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_SERVER")
	smtpPort := 587
	userEmail := os.Getenv("USER_EMAIL")

	m := gomail.NewMessage()
	m.SetHeader("From", senderName)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	d := gomail.NewDialer(smtpHost, smtpPort, senderName, senderPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending verification code: %v", err)
	}

	return nil
}
