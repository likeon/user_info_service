package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ExternalID  string    `gorm:"primaryKey" json:"external_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func main() {
	// gin router
	r := gin.Default()

	// db
	db, err := setupDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// migrate schema on strartup
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// routes
	r.POST("/save", func(c *gin.Context) {
		var input User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body: " + err.Error(),
			})
			return
		}

		// validate external_id as uuid
		if _, err := uuid.Parse(input.ExternalID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "external_id must be a valid UUID",
			})
			return
		}

		// db insert
		if err := db.Create(&input).Error; err != nil {
			// check for unique constraint error
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				c.JSON(http.StatusConflict, gin.H{
					"error": "external_id must be unique",
				})
				return
			}
			// generic error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save user: " + err.Error(),
			})
			return
		}

		// return http 201
		c.JSON(http.StatusCreated, input)
	})

	r.GET("/:external_id", func(c *gin.Context) {
		externalID := c.Param("external_id")

		// validate external_id as uuid
		if _, err := uuid.Parse(externalID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid external_id parameter",
			})
			return
		}

		var user User
		// db lookup
		result := db.First(&user, "external_id = ?", externalID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Record not found",
			})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// gorm db instance
func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
