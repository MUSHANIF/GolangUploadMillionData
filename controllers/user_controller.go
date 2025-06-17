package controllers

import (
	"encoding/csv"
	"net/http"
	"golangReact/database"
	"golangReact/models"
	"io"
	"golangReact/structs"
	"golangReact/helpers"
	"github.com/gin-gonic/gin"
	"fmt"
	"log"
	"time"
	"math/rand"
	"gorm.io/gorm"
)
var db *gorm.DB
func FindUsers(c *gin.Context) {

	// Inisialisasi slice untuk menampung data user
	var users []models.User

	// Ambil data user dari database
	database.DB.Find(&users)

	// Kirimkan response sukses dengan data user
	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Table user success",
		Data:    users,
	})
}

func CreateUser(c *gin.Context){
	req := structs.UserCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	user := models.User{
		Name:  req.Name,
		Username: req.Username,
		Email: req.Email,
		Password: helpers.HashPassword(req.Password),	
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}


	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "User created successfully",
		Data:  structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
func InsertMillionUsers(c *gin.Context) {
    batchSize := 10000
    total := 10000000
    var users []models.User

    start := time.Now()

    for i := 0; i < total; i++ {
        user := models.User{
            Name:     randomString(8),
            Username: fmt.Sprintf("user%d", i),
            Email:    fmt.Sprintf("user%d@example.com", i),
            Password: randomString(12),
			// Password: helpers.HashPassword("mypassword"), 
        }
        users = append(users, user)

        if len(users) >= batchSize {
            if err := database.DB.CreateInBatches(users, batchSize).Error; err != nil {
                log.Printf("Error inserting batch %d: %v", i/batchSize, err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            users = users[:0] 
            log.Printf("Inserted batch %d", i/batchSize)
        }
    }

    
    if len(users) > 0 {
        database.DB.CreateInBatches(users, batchSize)
    }

    duration := time.Since(start)
    c.JSON(http.StatusOK, gin.H{"message": "Inserted 10 million users", "duration": duration.String()})
}

func randomString(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func ImportUsersFromCSV(c *gin.Context) {
	    start := time.Now()
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	
	csvFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	
	_, err = reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read CSV header"})
		return
	}

	var users []models.User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading CSV rows"})
			return
		}

		user := models.User{
			Name:     record[0],
			Username: record[1],
			Email:    record[2],
			Password: record[3], 
		}

		users = append(users, user)
	}

	
	if err := database.DB.CreateInBatches(users, 1000).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	duration := time.Since(start)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Imported %d users", len(users)), "duration": duration.String()})
}

func FindUserById(c *gin.Context) {
	id := c.Param("id")

	user := models.User{}
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, structs.ErrorResponse{
				Success: false,
				Message: "User not found",
				Errors:  helpers.TranslateErrorMessage(err),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data: structs.UserResponse{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})

}
