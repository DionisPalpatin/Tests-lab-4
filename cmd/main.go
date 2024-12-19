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

	ctx.JSON(http.StatusOK, gin.H{"message": "Enter the TOTP code from your app"})
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

	if request.Email != userEmail || !totp.Validate(request.Code, totpSecret) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid TOTP code"})
		return
	}

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

	ctx.JSON(http.StatusOK, gin.H{"message": "Enter the TOTP code from your app"})
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

	if request.Email != userEmail || !totp.Validate(request.Code, totpSecret) {
		fmt.Print(request.Email, userEmail)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid TOTP code"})
		return
	}

	userPassword = request.NewPassword
	ctx.JSON(http.StatusOK, gin.H{"message": "Password is changed successfully"})
}
