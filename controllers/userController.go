package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"kharchasathi.com/database"
	helper "kharchasathi.com/helpers"
	"kharchasathi.com/models"
)

var userCollection = database.OpenCollection(database.Client, "users")

func GetUsers() gin.HandlerFunc {

	return func(c *gin.Context) {
		if !helper.IsAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this resource"})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get all users
		var users []models.User
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var user models.User
			cursor.Decode(&user)
			users = append(users, user)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUserById() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Check if the user is authenticated
		if !helper.IsAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this resource"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get the id f0rom the request params
		id := c.Param("id")

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Get the user with the given id
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SetUserPhoneNumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := helper.GetIdFromAccessToken(c)

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Create a new User object and bind the request body to it
		var user models.User
		err = c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if any other user has the same phone number
		var existingUser models.User
		err = userCollection.FindOne(ctx, bson.M{"phonenumber": user.PhoneNumber}).Decode(&existingUser)
		if err == nil && existingUser.ID != objId {
			// Remove the phone number from the existing user
			_, err = userCollection.UpdateOne(ctx, bson.M{"_id": existingUser.ID}, bson.M{"$unset": bson.M{"phonenumber": 1}})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Update the user with the given id
		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"phonenumber": user.PhoneNumber}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Phone number updated successfully"})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is authenticated
		if !helper.IsAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this resource"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get the id from the request params
		id := c.Param("id")

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Create a new User object and bind the request body to it
		var user models.User
		err = c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the user with the given id
		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": user})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is authenticated
		if !helper.IsAuthenticated(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to access this resource"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get the id from the request params
		id := c.Param("id")

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Delete the user with the given id
		result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func ChangeUserPassword() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := helper.GetIdFromAccessToken(c)
		var changePassword models.ChangePassword
		if err := c.BindJSON(&changePassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Get the user with the given id
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Verify old password
		if !helper.VerifyPassword(user.Password, changePassword.OldPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
			return
		}

		// Update the user with the given id
		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"password": helper.HashPassword(changePassword.NewPassword)}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}

func SignUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Create a new User object and bind the request body to it
		var user models.User
		if err := c.BindJSON(&user); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the email already exists
		existingUser := &models.User{}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}

		user.ID = primitive.NewObjectID()
		user.Password = helper.HashPassword(user.Password)

		_, err = userCollection.InsertOne(ctx, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Create a new LoginRequest object and bind the request body to it
		var userLogin models.UserLogin
		if err := c.BindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user with the given email exists
		existingUser := &models.User{}
		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(existingUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Check if the password is correct
		if !helper.VerifyPassword(existingUser.Password, userLogin.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate access token and refresh token and return
		accessToken, err := helper.GenerateAccessToken(existingUser.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		refreshToken, err := helper.GenerateRefreshToken(existingUser.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "User logged in successfully",
			"firstName":    existingUser.FirstName,
			"lastName":     existingUser.LastName,
			"email":        existingUser.Email,
			"phoneNumber":  existingUser.PhoneNumber,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func ChangeInitialCurrency() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := helper.GetIdFromAccessToken(c)
		var changeInitialCurrency models.ChangeInitialCurrency
		if err := c.BindJSON(&changeInitialCurrency); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Get the user with the given id
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		// Update the user with the given id
		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"initialcurrency": changeInitialCurrency.InitialCurrency}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Initial Currency updated successfully"})
	}
}
func ChangeInitialBalance() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := helper.GetIdFromAccessToken(c)
		var changeInitialBalance models.ChangeInitialBalance
		if err := c.BindJSON(&changeInitialBalance); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Get the user with the given id
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		// Update the user with the given id
		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"balance": changeInitialBalance.InitialBalance}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Initial balance updated successfully"})
	}
}
func AddIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		id := helper.GetIdFromAccessToken(c)
		var income models.Income
		if err := c.BindJSON(&income); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert the id to an ObjectId
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Get the user with the given id
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Add the income to the user's balance
		user.Balance += income.Amount

		// Update the user in the database
		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"balance": user.Balance}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Income added successfully"})
	}
}

func LoginWithGoogle() gin.HandlerFunc {
	return func(c *gin.Context) {

		type GoogleLoginRequest struct {
			TokenID string `json:"token_id"`
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Create a new GoogleLoginRequest object and bind the request body to it
		var googleLoginRequest GoogleLoginRequest
		if err := c.BindJSON(&googleLoginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify the token with Google and retrieve user profile information
		email, err := helper.VerifyGoogleToken(ctx, googleLoginRequest.TokenID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := helper.GetUserByEmail(email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Generate access token and refresh token and return
		accessToken, err := helper.GenerateAccessToken(user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		refreshToken, err := helper.GenerateRefreshToken(user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "User logged in successfully",
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"phoneNumber":  user.PhoneNumber,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}
func GetUser(id primitive.ObjectID) (models.User, error) {
	// Get users collection from database

	// Define filter to retrieve user by id
	filter := bson.M{"_id": id}

	// Retrieve user from database
	var user models.User
	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return nil user and error if user not found
			return models.User{}, fmt.Errorf("user not found")
		} else {
			// Return nil user and error for other errors
			return models.User{}, err
		}
	}

	// Return retrieved user
	return user, nil
}
