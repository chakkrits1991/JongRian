package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	// แก้ไขข้อมูลเชื่อมต่อตรงนี้ให้ตรงกับเครื่องคุณ
	dsn := "host=localhost port=5432 user=postgres password=1234 dbname=postgres sslmode=disable"

	var err error
	DB, err = sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalln("เปิดการเชื่อมต่อ DB ไม่ได้:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalln("เชื่อมต่อ DB ไม่สำเร็จ (Ping fail):", err)
	}

	fmt.Println("✅ เชื่อมต่อ PostgreSQL สำเร็จ!")
}
