package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"net/http"
)

func (h *Handler) statusHandler(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (h *Handler) loginHandler(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Password != h.UserPassword || request.Email != h.UserEmail {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Enter the TOTP code from your app"})
}

func (h *Handler) verifyLoginHandler(ctx *gin.Context) {
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Email != h.UserEmail || !totp.Validate(request.Code, h.TotpSecret) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid TOTP code"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Auth is made successfully"})
}

func (h *Handler) resetPasswordHandler(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Password != h.UserPassword || request.Email != h.UserEmail {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Enter the TOTP code from your app"})
}

func (h *Handler) verifyResetHandler(ctx *gin.Context) {
	var request struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Email != h.UserEmail || !totp.Validate(request.Code, h.UserEmail) {
		fmt.Print(request.Email, h.UserEmail)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid TOTP code"})
		return
	}

	h.UserPassword = request.NewPassword
	ctx.JSON(http.StatusOK, gin.H{"message": "Password is changed successfully"})
}
