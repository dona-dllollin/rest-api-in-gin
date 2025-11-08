package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest-api-in-gin/internal/database"
	"rest-api-in-gin/internal/helper"
	"strconv"
	"time"

	"rest-api-in-gin/internal/env"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (app *application) login(c *gin.Context) {
	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		helper.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if existingUser == nil {
		helper.JSONError(c, http.StatusUnauthorized, "Invalid email or password")
		fmt.Println(err.Error())
		return
	}
	if err != nil {
		helper.JSONError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		helper.JSONError(c, http.StatusUnauthorized, "Invalid email or password")
		fmt.Println(*existingUser)
		return
	}

	tokenString, err := helper.GenerateAccessToken(existingUser.Id)
	refreshToken, _ := helper.GenerateRefreshToken(existingUser.Id)
	if err != nil {
		helper.JSONError(c, http.StatusInternalServerError, "gagal membuat token")
	}

	// simpan refresh token kdi redis
	app.redis.Set(context.Background(), "refresh:"+strconv.Itoa(existingUser.Id), refreshToken, 7*24*time.Hour)

	//kirim refresh token via httpOnly Cookie
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "", true, true)

	c.JSON(http.StatusOK, loginResponse{Token: tokenString})

}

func (app *application) registerUser(c *gin.Context) {
	var register registerRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPasswod, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Something wnet wrong"})
		return
	}

	register.Password = string(hashedPasswod)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Name:     register.Name,
	}
	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		log.Println(user)
		return
	}
	c.JSON(http.StatusCreated, user)
}

// handler untuk refresh access token

func (app *application) Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil || cookie == "" {
		helper.JSONError(c, http.StatusUnauthorized, "No refresh token")
		return
	}
	// parse refresh token
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(env.GetEnvString("REFRESH_SECRET", "some-iam-fresh-2718271")), nil
	})
	if err != nil || !token.Valid {
		helper.JSONError(c, http.StatusUnauthorized, "Invalid Refresh token")
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userId"].(float64)

	// Ambil token dari redis untuk validasi
	stored, err := app.redis.Get(context.Background(), "refresh:"+strconv.Itoa(int(userID))).Result()
	if err != nil || stored != cookie {
		helper.JSONError(c, http.StatusUnauthorized, "Invalid Refresh token")
		return
	}

	// Buat Access Token Baru
	user_id := int(userID)
	newAccessToken, _ := helper.GenerateAccessToken(user_id)

	c.JSON(http.StatusOK, loginResponse{Token: newAccessToken})

}
