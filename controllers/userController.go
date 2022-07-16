package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/VarunSAthreya/go-jwt/database"
	"github.com/VarunSAthreya/go-jwt/helpers"
	"github.com/VarunSAthreya/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	if err != nil {
		return false, "Email or Password is incorrect."
	}

	return true, ""
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking Email."})

			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists."})

			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking Phone number."})

			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone Number already exists."})

			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, tokenErr := helpers.GenerateTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)

		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while generating tokens."})

			return
		}

		user.Token = &token
		user.Refresh_token = &refreshToken

		password, passError := HashPassword(*user.Password)
		if passError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while hashing password."})

			return
		}

		user.Password = &password

		resultInsertNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while inserting user."})

			return
		}

		c.JSON(http.StatusOK, resultInsertNumber)
	}
}

func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user = models.User{}
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or Password is incorrect."})

			return
		}

		passwordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})

			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error user Email not found."})

			return
		}

		token, refreshToken, tokenErr := helpers.GenerateTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)

		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while generating tokens."})

			return
		}

		helpers.UpdateTokens(foundUser.User_id, token, refreshToken)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}

		c.JSON(http.StatusOK, user)
	}
}
