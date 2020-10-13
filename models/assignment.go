package models

import (
	"gorone/db"
	"time"

	"gorm.io/gorm"
)

// Assignment Requestに紐づくworker
type Assignment struct {
	RequestID int
	PublicIP  string
	gorm.Model
	TerminatedAt *time.Time
}

// NewAssignment 紐付けの作成
func NewAssignment(requestID int, ip string) (*Assignment, error) {

	assignment := Assignment{RequestID: requestID, PublicIP: ip}
	db := db.DbManager()
	if err := db.Create(&assignment).Error; err != nil {
		return nil, err
	}

	return &assignment, nil
}

// Terminate workerを終了するときに発行
func (a Assignment) Terminate() error {
	db := db.DbManager()
	now := time.Now()
	a.TerminatedAt = &now

	if err := db.Save(&a).Error; err != nil {
		return err
	}

	return nil
}
