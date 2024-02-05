package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SanjaySinghRajpoot/newsFeed/config"
	"github.com/SanjaySinghRajpoot/newsFeed/models"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/formatError"
	helpers "github.com/SanjaySinghRajpoot/newsFeed/utils/helper"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/pagination"
	"github.com/SanjaySinghRajpoot/newsFeed/utils/validations.go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Signup function is used to create a user or signup a user
func Signup(c *gin.Context) {
	// Get the name, email and password from request
	var userInput struct {
		Name     string `json:"name" binding:"required,min=2,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errs),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Email unique validation
	if validations.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"Email": "The email is already exist!",
			},
		})
		return
	}

	// Hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	user := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: string(hashPassword),
	}

	// Create the user
	result := config.DB.Create(&user)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Login function is used to log in a user
func Login(c *gin.Context) {
	// Get the email and password from the request
	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if c.ShouldBindJSON(&userInput) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// GetTokenString, er := redis.GetRedisData(userInput.Email)
	// if er != nil {
	// 	fmt.Printf("Failed to Get the Redis Cache, Setting the Cache: %s", er)
	// }

	// if GetTokenString != "" {

	// 	fmt.Println("redis cache working for login")

	// 	c.SetSameSite(http.SameSiteLaxMode)
	// 	c.SetCookie("Authorization", GetTokenString, 3600*24*30, "", "", false, true)
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "User login successful",
	// 	})

	// 	return
	// }

	// Find the user by email
	var user models.User
	config.DB.First(&user, "email = ?", userInput.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Compare the password with user hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign in and get the complete encoded token as a string using the .env secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// set the redis cache here
	// msg, error := redis.SetRedisData(userInput.Email, tokenString)

	// if error != nil {

	// 	fmt.Printf("Failed to Set the Redis Cache: %s", msg)

	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": error.Error(),
	// 	})

	// 	return
	// }

	// Set expiry time and send the token back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "User login successful",
	})
}

// Logout function is used to log out a user
func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("Authorization", "", 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"successMessage": "Logout successful",
	})
}

// GetUsers function is used to get users list
func GetUsers(c *gin.Context) {
	// Get all the users
	var users []models.User

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, _ := strconv.Atoi(perPageStr)

	result, err := pagination.Paginate(config.DB, page, perPage, nil, &users)
	if err != nil {
		formatError.InternalServerError(c, err)
		return
	}

	// Return the users
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

// // EditUser function is used to find a user by id
func GetUser(c *gin.Context) {
	// Get the id from url
	id := c.Param("id")

	// Find the user
	var user models.User
	result := config.DB.First(&user, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Return the user
	c.JSON(http.StatusOK, gin.H{
		"result": user,
	})
}

// UpdateUser function is used to update a user
func UpdateUser(c *gin.Context) {
	// Get the id from url
	id := c.Param("id")

	UserID := helpers.GetAuthUser(c).ID

	userIDStr := fmt.Sprintf("%d", UserID)

	if id != userIDStr {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User can only update their own information",
		})
		return
	}

	// Get the name, email and password from request
	var userInput struct {
		Name     string `json:"name" binding:"required,min=2,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": validations.FormatValidationErrors(errs),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find the user by id
	var user models.User
	result := config.DB.First(&user, id)

	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Email unique validation
	if user.Email != userInput.Email && validations.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"Email": "The email is already exist!",
			},
		})
		return
	}

	// Prepare data to update
	updateUser := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	// Update the user
	result = config.DB.Model(&user).Updates(&updateUser)

	if result.Error != nil {
		formatError.InternalServerError(c, result.Error)
		return
	}

	// Return the user
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// DeleteUser function is used to delete a user by id
func DeleteUser(c *gin.Context) {
	// Get the id from the url
	id := c.Param("id")
	var user models.User

	result := config.DB.First(&user, id)
	if err := result.Error; err != nil {
		formatError.RecordNotFound(c, err)
		return
	}

	// Delete the user
	config.DB.Delete(&user)

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "The user has been deleted successfully",
	})
}
