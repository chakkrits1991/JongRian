package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// User - ข้อมูลผู้ใช้งาน
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	FullName  string    `json:"full_name" db:"full_name"`
	Role      string    `json:"role" db:"role"` // 'tutor' หรือ 'student'
	LineID    string    `json:"line_id,omitempty" db:"line_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Course - วิชาที่เปิดสอน
type Course struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	TutorID      uuid.UUID       `json:"tutor_id" db:"tutor_id"`
	Title        string          `json:"title" db:"title"`
	Description  string          `json:"description,omitempty" db:"description"`
	PricePerHour decimal.Decimal `json:"price_per_hour" db:"price_per_hour"`
}

// AvailabilitySlot - ช่วงเวลาที่ว่างสำหรับจอง
type AvailabilitySlot struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TutorID   uuid.UUID `json:"tutor_id" db:"tutor_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
	Status    string    `json:"status" db:"status"` // 'available', 'booked', 'pending'
}

// Booking - ข้อมูลการจองและการชำระเงิน
type Booking struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CourseID       uuid.UUID `json:"course_id" db:"course_id"`
	StudentID      uuid.UUID `json:"student_id" db:"student_id"`
	SlotID         uuid.UUID `json:"slot_id" db:"slot_id"`
	Status         string    `json:"status" db:"status"` // 'pending_payment', 'confirmed', 'cancelled'
	PaymentSlipURL string    `json:"payment_slip_url,omitempty" db:"payment_slip_url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
