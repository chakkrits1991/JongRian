package repository

import (
	"context"
	"database/sql"
	"errors"
	"jongrian/models"

	"github.com/google/uuid"
)

// CreateBooking ทำหน้าที่จัดการการจองในรูปแบบ Transaction
func CreateBooking(db *sql.DB, studentID uuid.UUID, courseID uuid.UUID, slotID uuid.UUID) (*models.Booking, error) {
	// 1. เริ่มต้น Transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// 2. ตรวจสอบสถานะ Slot และ Lock แถวนั้นไว้ (SELECT FOR UPDATE)
	// เพื่อไม่ให้ Transaction อื่นมาแก้ไขในระหว่างที่เรากำลังทำงาน
	var status string
	err = tx.QueryRow("SELECT status FROM availability_slots WHERE id = $1 FOR UPDATE", slotID).Scan(&status)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if status != "available" {
		tx.Rollback()
		return nil, errors.New("ขออภัยครับ ช่วงเวลานี้ถูกจองไปแล้ว")
	}

	// 3. อัปเดตสถานะ Slot เป็น 'pending' (รอชำระเงิน)
	_, err = tx.Exec("UPDATE availability_slots SET status = 'pending' WHERE id = $1", slotID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4. สร้าง Record การจองใหม่
	newBooking := models.Booking{
		ID:        uuid.New(),
		CourseID:  courseID,
		StudentID: studentID,
		SlotID:    slotID,
		Status:    "pending_payment",
	}

	query := `INSERT INTO bookings (id, course_id, student_id, slot_id, status) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, newBooking.ID, newBooking.CourseID, newBooking.StudentID, newBooking.SlotID, newBooking.Status)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 5. Commit ข้อมูลทั้งหมดลง Database
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &newBooking, nil
}
