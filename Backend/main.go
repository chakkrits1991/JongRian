package main

import (
	"jongrian/database"
	"jongrian/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware สำหรับแก้ปัญหา CORS ให้ Frontend คุยกับ Backend ได้
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	database.ConnectDB()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// 1. ดึงรายการเวลาทั้งหมด (GET)
	r.GET("/api/slots", func(c *gin.Context) {
		var slots []models.AvailabilitySlot
		err := database.DB.Select(&slots, "SELECT * FROM availability_slots ORDER BY start_time ASC")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, slots)
	})

	// 2. เพิ่มเวลาว่างใหม่ (POST)
	r.POST("/api/slots", func(c *gin.Context) {
		var newSlot models.AvailabilitySlot
		if err := c.ShouldBindJSON(&newSlot); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
			return
		}

		query := `INSERT INTO availability_slots (start_time, end_time, status) 
				  VALUES (:start_time, :end_time, 'available')`

		_, err := database.DB.NamedExec(query, newSlot)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกไม่สำเร็จ: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "เพิ่มเวลาว่างสำเร็จ!"})
	})

	// 3. ระบบกดจองเวลา (PATCH)
	r.PATCH("/api/slots/:id/book", func(c *gin.Context) {
		id := c.Param("id")
		query := `UPDATE availability_slots SET status = 'booked' WHERE id = $1 AND status = 'available'`

		result, err := database.DB.Exec(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "จองไม่สำเร็จ"})
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "เวลานี้ถูกจองไปแล้วหรือไม่มีอยู่จริง"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "จองสำเร็จ!"})
	})

	r.Run(":8080")
}
