package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"log"
	"net/http"
	"os"
)

var totpSecret string
var userPassword string
var userEmail string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	userPassword = os.Getenv("USER_PASSWORD")
	userEmail = os.Getenv("USER_EMAIL")

	fmt.Print("user email: ", userEmail)

	initTOTP()

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

func initTOTP() {
	totpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "BDDTests",
		AccountName: userEmail,
	})

	if err != nil {
		log.Fatalf("Error generating TOTP secret: %v", err)
	}

	if os.Getenv("THIS_IS_TEST") != "" {
		totpSecret = os.Getenv("TOTP_SECRET")
	} else {
		totpSecret = totpKey.Secret()
	}
}
