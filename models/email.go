package models

import (
	"log"

	"gorm.io/gorm"
)

const (
	EmailTableName = "verificationCode"
)

type Emailverify struct {
	Email string `gorm:"column:email" json:"email"`
	Code  string `gorm:"column:code" json:"code"`
}
type EmailAddress struct {
	Email string `json:"email"`
}
type Code struct {
	Code string `json:"code"`
}

func (e *Emailverify) TableName() string {
	return EmailTableName
}

// 定时清除生成的验证码以保证验证码有效时间为五分钟
func CreateEvent(db *gorm.DB) {
	eventSQL := `
        CREATE EVENT IF NOT EXISTS delete_verificationCode
        ON SCHEDULE EVERY 5 MINUTE
        DO
          DELETE FROM verificationCode WHERE code < NOW() - INTERVAL 5 MINUTE;
    `
	if err := db.Exec(eventSQL).Error; err != nil {
		log.Println("Error creating event:", err)
	} else {
		log.Println("Event created successfully")
	}
}
